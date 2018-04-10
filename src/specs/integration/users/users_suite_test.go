package users

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"

	"testing"

	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Test Suite")
}

const defaultAdminPassword = "REPLACE-ME"

var boshDeployment = os.Getenv("ENGINE_DEPLOYMENT")

var _ = BeforeSuite(func() {
	requiredEnvs := []string{
		"BOSH_CA_CERT",
		"BOSH_CLIENT",
		"BOSH_CLIENT_SECRET",
		"ENGINE_DEPLOYMENT",
		"BOSH_ENVIRONMENT",
	}

	CheckForRequiredEnvVars(requiredEnvs)
})

func showGrants(username, password string) string {
	sql := "SHOW GRANTS"
	session := MustSucceed(ExecuteMysqlQuery(boshDeployment, "0", username, password, sql))
	return string(session.Out.Contents())
}

func showTable(username, password, schemaName, tableName string) string {
	sql := fmt.Sprintf(
		`SELECT TABLE_SCHEMA, TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'`,
		schemaName, tableName)
	session := MustSucceed(ExecuteMysqlQuery(boshDeployment, "0", username, password, sql))
	return string(session.Out.Contents())
}

func queryDatabaseAsUser(username, password string) *gexec.Session {
	return ExecuteMysqlQuery(boshDeployment, "0", username, password, "SELECT 1")
}
