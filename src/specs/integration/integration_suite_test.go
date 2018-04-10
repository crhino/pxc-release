package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	requiredEnvs := []string{
		"BOSH_CA_CERT",
		"BOSH_CLIENT",
		"BOSH_CLIENT_SECRET",
		"ENGINE_DEPLOYMENT",
		"BOSH_ENVIRONMENT",
	}

	testhelpers.CheckForRequiredEnvVars(requiredEnvs)
})
