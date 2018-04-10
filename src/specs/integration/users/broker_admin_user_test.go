package users

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"
)

var _ = Describe("Broker Admin User", func() {
	const defaultPassword = "REPLACE-BROKER-ADMIN-PASSWORD"

	It("broker-admin user has the ability to grant privileges to other users", func() {
		session := testhelpers.ExecuteMysqlQuery(boshDeployment, "0", "broker-admin", defaultPassword, "SHOW GRANTS")
		Expect(session).To(gexec.Exit(0))
		output := string(session.Out.Contents())

		By("being granted all privileges on all databases", func() {
			Expect(output).Should(ContainSubstring("GRANT ALL PRIVILEGES ON `%`.* TO 'broker-admin'@'%' WITH GRANT OPTION"))
		})
		By("being granted the privilege to create new users", func() {
			Expect(output).Should(ContainSubstring("GRANT CREATE USER ON *.* TO 'broker-admin'@'%'"))
		})
		By("being disallowed from directly modifying the mysql system database", func() {
			Expect(output).Should(ContainSubstring("GRANT SHOW VIEW ON `mysql`.* TO 'broker-admin'@'%' WITH GRANT OPTION"))
		})
	})
})
