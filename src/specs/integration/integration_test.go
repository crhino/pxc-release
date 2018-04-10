package integration_test

import (
	"os"
	"time"

	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var boshDeployment = os.Getenv("ENGINE_DEPLOYMENT")

var _ = Describe("Engine", func() {
	Context("drain script", func() {
		It("succeeds when mysql is already stopped", func() {
			Expect(
				testhelpers.ExecuteBosh([]string{
					"-d", boshDeployment,
					"ssh",
					"-c",
					"sudo /var/vcap/bosh/bin/monit unmonitor mysql && sudo /var/vcap/jobs/mysql/bin/mysql_ctl stop",
				}, 2*time.Minute),
			).To(gexec.Exit(0))

			Expect(

				testhelpers.ExecuteBosh([]string{
					"-d", boshDeployment,
					"-n",
					"restart", "mysql",
				}, 2*time.Minute),
			).To(gexec.Exit(0))
		})
	})
})
