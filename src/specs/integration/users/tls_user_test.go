package users

import (
	"time"

	. "specs/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"
)

var _ = Describe("Enforcing client TLS", func() {
	Context("when enforce_client_tls is enabled", func() {
		BeforeEach(func() {
			enableEnforceClientTLSOpsFileContents := `
				- type: replace
				  path: /properties/enforce_client_tls?
				  value: true
			`

			RedeployWithOpsFile(boshDeployment, testhelpers.ProperYaml(enableEnforceClientTLSOpsFileContents))
		})

		AfterEach(func() {
			disableEnforceClientTLSOpsFileContents := `
				- type: replace
				  path: /properties/enforce_client_tls?
				  value: false
			`

			RedeployWithOpsFile(boshDeployment, testhelpers.ProperYaml(disableEnforceClientTLSOpsFileContents))

		})

		It("cannot connect through an insecure channel", func() {
			mysqlCmd := `mysql --defaults-file=/var/vcap/jobs/mysql/config/mylogin.cnf --ssl-mode=disabled --host=127.0.0.1`
			args := []string{
				"-d", boshDeployment,
				"ssh",
				"mysql/0",
				"-c", mysqlCmd,
			}
			session := testhelpers.ExecuteBosh(args, 2*time.Minute)
			Eventually(session, "2m").Should(gexec.Exit(1))
			Expect(session).To(gbytes.Say(`ERROR 3159 \(HY000\): Connections using insecure transport are prohibited while --require_secure_transport=ON`))
		})
	})

})
