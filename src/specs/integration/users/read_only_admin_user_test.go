package users

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "specs/test_helpers"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"
)

var _ = Describe("Read-Only Admin User", func() {
	const defaultPassword = "REPLACE-ME-READ-ONLY"

	Context("when the read_only_admin_password is provided", func() {
		BeforeEach(func() {
			addReadOnlyAdminUser := `
				- type: replace
				  path: /properties/read_only_admin_password?
				  value: REPLACE-ME-READ-ONLY
			`

			RedeployWithOpsFile(boshDeployment, testhelpers.ProperYaml(addReadOnlyAdminUser))
		})

		It("can issue a query as the roadmin user", func() {
			session := queryDatabaseAsUser("roadmin", defaultPassword)
			Expect(session).To(gexec.Exit(0))
		})

		It("grants select and process to the read only admin user", func() {
			output := showGrants("roadmin", defaultPassword)
			Expect(output).To(ContainSubstring(`GRANT SELECT, PROCESS ON *.* TO 'roadmin'@'%'`))
		})
	})

	Context("when the read_only_admin_password has been removed", func() {
		BeforeEach(func() {
			removeReadOnlyAdminUserOpsFileContents := `
				- type: remove
				  path: /properties/read_only_admin_password?
		    `

			RedeployWithOpsFile(boshDeployment, testhelpers.ProperYaml(removeReadOnlyAdminUserOpsFileContents))
		})

		It("cannot connect as the roadmin user", func() {
			session := queryDatabaseAsUser("roadmin", defaultPassword)
			Expect(session).NotTo(gexec.Exit(0))
			Expect(session).To(gbytes.Say("Access denied for user 'roadmin'"))
		})
	})

})
