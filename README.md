# pxc-release

Alpha Percona Xtradb Cluster release **Not ready for production use**

pxc-release is a BOSH release of MySQL Galera that can be used as a backing store for Cloudfoundry. The Galera Cluster Provider is [Percona Xtradb Cluster](https://www.percona.com/software/mysql-database/percona-xtradb-cluster). 
This release is intended as a drop-in replacement for [cf-mysql-release](https://github.com/cloudfoundry/cf-mysql-release).

<a name='components'></a>
# Components

<a name='mysql-server'></a>
## MySQL Server

<a name='proxy'></a>
## Proxy

Traffic to the MySQL cluster is routed through one or more proxy nodes. The current proxy implementation is [Switchboard](https://github.com/cloudfoundry-incubator/switchboard). This proxy acts as an intermediary between the client and the MySQL server, providing failover between MySQL nodes. The number of nodes is configured by the proxy job instance count in the deployment manifest.

**NOTE:** If the number of proxy nodes is set to zero, apps will be bound to the IP address of the first MySQL node in the cluster. If that IP address should change for any reason (e.g. loss of a VM) or a proxy was subsequently added, one would need to re-bind all apps to the IP address of the new node.

For more details see the [proxy documentation](/docs/proxy.md).

<a name='deploying'></a>
# Deploying

## Deployment Components and Topology

The typical topology is 2 proxy nodes and 3 mysql-clustered nodes. The proxies can be separate vms or co-located with the mysql-clustered nodes.

### Database nodes

The number of mysql nodes should always be odd, with a minimum count of three, to avoid [split-brain](http://en.wikipedia.org/wiki/Split-brain\_\(computing\)).
When a failed node comes back online, it will automatically rejoin the cluster and sync data from one of the healthy nodes.

### Proxy nodes

Two proxy instances are recommended. The second proxy is intended to be used in a failover capacity. You can also choose to place a load balancer in front of both proxies, or use [BOSH DNS](https://bosh.io/docs/dns.html) to send traffic to both proxies.

In the event the first proxy fails, the second proxy will still be able to route requests to the mysql nodes.

The proxies both will route traffic to the lowest-indexed healthy galera node, according to the galera index (not bosh index).
<a name='deploying-new-deployments'></a>
## Deploying new deployments
<a name='deploying-with-cf-deployment'></a>
### Deploying CF with pxc-release
Use the [cf-deployment manifests](https://github.com/cloudfoundry/cf-deployment) with the `experimental/use-pxc.yml` ops file. (Currently in a PR here: https://github.com/cloudfoundry/cf-deployment/pull/453/files soon to be widely available)

<a name='deploying-with-non-cf-deployments'></a>
### Using PXC release with other deployments
1. Get the latest pxc bosh release from [bosh.io](http://bosh.io/releases/github.com/cloudfoundry-incubator/pxc-release)
2. Add the release from bosh.io to your manifest
3. Configure the properties from the job spec sections for the [mysql-clustered](jobs/mysql-clustered/spec) and the [proxy](jobs/mysql-clustered-spec) jobs. You can use the manifest and ops-files in cf-deployment as a guide to configuring these properties. See [Deploying CF with pxc-release](#deploying-with-cf-deployment)

<a name='migrating-with-cf-deployment'></a>
## Migrating from cf-mysql-release

Requirements:

[cf-mysql-release](https://github.com/cloudfoundry/cf-mysql-release/) v36.12.0 or greater

<a name='migrating-with-cf-deployment'></a>
### Migrating CF with pxc-release
Use the [cf-deployment manifests](https://github.com/cloudfoundry/cf-deployment) with the `experimental/migrate-cf-mysql-to-pxc.yml` ops file. (Currently in a PR here: https://github.com/cloudfoundry/cf-deployment/pull/453/files soon to be widely available) It is advisable to take a backup first.

The ops file will trigger the same migration procedure described in [Using PXC release with other deployments](#migrating-with-non-cf-deployments)

After migrating, use the [Deploying CF with pxc-release](#deploying-with-cf-deployment) docs for your next deploy.

<a name='migrating-with-non-cf-deployments'></a>
### Using PXC release with other deployments

1. Make backups according to your usual backup procedure.
1. Get the latest pxc bosh release from [bosh.io](http://bosh.io/releases/github.com/cloudfoundry-incubator/pxc-release)
2. Add the release to your manifest
2. ⚠️ **Scale down to 1 node and ensure the persistent disk has enough free space to double the size of the mysql data.**
3. Make the following changes to your bosh manifest:
   * Add the `mysql-clustered` job from `pxc-release` to the instance group that has the `mysql` job from `cf-mysql-release`
   * Configure the `mysql-clustered` job with the same credentials and property values as the `mysql` job
   * To run the migration:
      * Set the `cf_mysql_enabled: false` property on the `mysql` job
      * Set the `pxc_enabled: true` property on `mysql-clustered` job
      * Switch the proxies to use the proxy job from `pxc-release` instead of `cf-mysql-release`
      * Deploy using BOSH

   * To prepare for the migration, but not run it immediately:
      * Set the `cf_mysql_enabled: true` property on the `mysql` job
      * Set the `pxc_enabled: false` property on `mysql-clustered` job
      * Deploy using BOSH
      * The MySQL will run as normal with only the `cf-mysql-release` running
      * In order to trigger the migration, redeploy with `cf_mysql_enabled: false` and `pxc_enabled: true`

   * ⚠️ **Do not enable both releases or disable both releases. Only enable one at a time.**
4. The migration is triggered by deploying with `cf_mysql_enabled: false` and `pxc_enabled: true`. The `pre-start` script for the `mysql-clustered` job in `pxc-release` starts both the Mariadb MySQL from the `cf-mysql-release` and the Percona MySQL from `pxc-release`. The migration dumps the MariaDB MySQL and loads that data into the Percona MySQL. This is done using pipes, so the dump is not written to disk, in order to reduce the use of disk space. The MariaDB MySQL is then stopped, leaving only the Percona MySQL running.
   * ⚠️ **MySQL DB will experience downtime during the migration**
5. After the migration, you can optionally clean up your deployment:
   * The migration will make a copy of the MySQL data on the persistent disk. To reduce disk usage, you can delete the old copy of the data in `/var/vcap/store/mysql` after you feel comfortable in the success of your migration. Do **NOT** delete the new copy of the data in `/var/vcap/store/mysql-clustered`.
   * Deploy only the `pxc-release` and not the `cf-mysql-release` in future deployments per [Deploying new deployments](#deploying-new-deployments), to free up disk space used by the `cf-mysql-release`.
   
6. Scale back up to the recommended 3 nodes, if desired.

<a name='contribution-guide'></a>
# Contribution Guide

The Cloud Foundry team uses GitHub and accepts contributions via
[pull request](https://help.github.com/articles/using-pull-requests).

## Contributor License Agreement

Follow these steps to make a contribution to any of our open source repositories:

1. Ensure that you have completed our CLA Agreement for
  [individuals](https://www.cloudfoundry.org/pdfs/CFF_Individual_CLA.pdf) or
  [corporations](https://www.cloudfoundry.org/pdfs/CFF_Corporate_CLA.pdf).

1. Set your name and email (these should match the information on your submitted CLA)

        git config --global user.name "Firstname Lastname"
        git config --global user.email "your_email@example.com"

## General Workflow

1. Fork the repository
1. Check out `develop` of pxc-release
1. Create a feature branch (`git checkout -b <my_new_branch>`)
1. Make changes on your branch
1. Deploy your changes using pxc as the database for cf-deployment to your dev environment and run [CF Acceptance Tests (CATS)](https://github.com/cloudfoundry/cf-acceptance-tests/)
1. Push to your fork (`git push origin <my_new_branch>`) and submit a pull request

We favor pull requests with very small, single commits with a single purpose.

Your pull request is much more likely to be accepted if it is small and focused with a clear message that conveys the intent of your change.
