package users

import (
	. "specs/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"
)

var _ = Describe("Metrics User", func() {
	const defaultPassword = "REPLACE-METRICS-PASSWORD"

	Context("when the mysql_metrics_password is provided", func() {
		BeforeEach(func() {
			addMetricsUserOpsFileContents := `
				- type: replace
				  path: /properties/mysql_metrics_password?
				  value: REPLACE-METRICS-PASSWORD
			`

			RedeployWithOpsFile(boshDeployment, testhelpers.ProperYaml(addMetricsUserOpsFileContents))
		})

		It("can issue a query as the mysql-metrics user", func() {
			session := queryDatabaseAsUser("mysql-metrics", defaultPassword)
			Expect(session).To(gexec.Exit(0))
		})
	})

	Context("when the mysql_metrics_password has been removed", func() {
		BeforeEach(func() {
			removeMetricsUserOpsFileContents := `
				- type: remove
				  path: /properties/mysql_metrics_password?
		    `

			RedeployWithOpsFile(boshDeployment, testhelpers.ProperYaml(removeMetricsUserOpsFileContents))
		})

		It("cannot connect as the mysql-metrics user", func() {
			session := queryDatabaseAsUser("mysql-metrics", defaultPassword)
			Expect(session).NotTo(gexec.Exit(0))
			Expect(session).To(gbytes.Say("Access denied for user 'mysql-metrics'"))
		})
	})
})
