package test_helpers

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"fmt"

	"github.com/pivotal-cf/dedicated-mysql-utils/testhelpers"
)

const boshPath = "/usr/local/bin/bosh"

// DisableResurrector disables the health monitor resurrector for deploymentName
func DisableResurrector(deploymentName string) {
	command := exec.Command(
		boshPath, "update-resurrection", "-d", deploymentName, "off")
	err := command.Run()
	var errStr string
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			errStr = string(e.Stderr)
		}
	}
	Expect(err).NotTo(HaveOccurred(), errStr)
}

// CloudCheck runs bosh cloud-check for deploymentName
func CloudCheck(deploymentName string) {
	command := exec.Command(
		boshPath, "-n", "cck", "-d", deploymentName, "-a")
	command.Stdout = GinkgoWriter
	command.Stderr = GinkgoWriter

	err := command.Run()
	Expect(err).NotTo(HaveOccurred())
}

func WaitForResurrector(deploymentName string, instanceIndex int) {
	Eventually(func() string {
		instance := testhelpers.GetMySQLInstancesSortedByIndex(deploymentName)[instanceIndex]
		return instance.ProcessState
	}, 15*time.Minute, time.Second).Should(Equal("running"))

	Eventually(func() *gexec.Session {
		return testhelpers.ExecuteBoshNoOutput([]string{"--tty", "-d", deploymentName, "tasks"}, time.Minute)
	}, 15*time.Minute, time.Second).Should(gbytes.Say("0 tasks"))
}

func RedeployWithOpsFile(boshDeployment string, opsFileContents []byte) string {
	args := []string{"-c", fmt.Sprintf("bosh -d %[1]s -n deploy <( bosh -d %[1]s manifest ) --ops-file <( echo '%[2]s' )", boshDeployment, opsFileContents)}
	output, err := exec.Command("bash", args...).CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))
	return string(output)
}

func ShowEvent(boshDeployment, schemaName, eventName string) string {
	sql := fmt.Sprintf("SHOW EVENTS IN %s WHERE name = '%s'\\G",
		schemaName, eventName)
	return testhelpers.ExecuteMysqlQueryAsAdmin(boshDeployment, "0", sql)
}

func DbSchemaExists(boshDeployment, schemaName string) bool {
	sql := fmt.Sprintf(`SELECT COUNT(*) = 1 FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'`, schemaName)
	result := testhelpers.ExecuteMysqlQueryAsAdmin(boshDeployment, "0", sql)
	return result == "1"
}

func DbTableExists(boshDeployment, schemaName, tableName string) bool {
	sql := fmt.Sprintf(`SELECT COUNT(*) = 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'`, schemaName, tableName)
	result := testhelpers.ExecuteMysqlQueryAsAdmin(boshDeployment, "0", sql)
	return result == "1"
}

const monitPath = "/var/vcap/bosh/bin/monit"

func ResetDbState(boshDeployment, instanceIndex string) {
	cmd := fmt.Sprintf(`sudo %[1]s unmonitor mysql && \
		sudo /var/vcap/jobs/mysql/bin/mysql_ctl stop && \
		sudo rm -rf /var/vcap/store/mysql/data && \
		sudo rm -f /var/vcap/sys/run/mysql/lf-state/leader.cnf && \
		sudo /var/vcap/jobs/mysql/bin/mysql_ctl start && \
		sudo %[1]s monitor mysql`, monitPath)
	session := testhelpers.ExecuteBosh([]string{"-d", boshDeployment, "ssh", "mysql/" + instanceIndex, "-c", cmd}, 2*time.Minute)
	Expect(session).To(gexec.Exit(0))
}

func GetInstanceCID(deploymentName string, instanceIndex int) string {
	return testhelpers.GetMySQLInstancesSortedByIndex(deploymentName)[instanceIndex].VmCid
}

func GetInstanceIP(deploymentName string, instanceIndex int) string {
	return testhelpers.GetMySQLInstancesSortedByIndex(deploymentName)[instanceIndex].IP
}
