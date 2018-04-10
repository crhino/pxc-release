package users

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"
)

var _ = Describe("Admin Database User", func() {

	It("can issue a query as admin", func() {
		session := queryDatabaseAsUser("admin", defaultAdminPassword)
		Expect(session).To(gexec.Exit(0))
	})

	It("cannot log in as root", func() {
		session := queryDatabaseAsUser("root", defaultAdminPassword)
		Expect(session).NotTo(gexec.Exit(0))
		Expect(session).To(gbytes.Say("Access denied for user 'root'"))
	})

	It("only can access from localhost", func() {
		output := testhelpers.ExecuteMysqlQueryAsAdmin(
			boshDeployment,
			"0",
			"SELECT host FROM mysql.user WHERE user = 'admin'")

		Expect(output).To(ContainSubstring("127.0.0.1"))
		Expect(output).To(ContainSubstring("localhost"))
	})
})
