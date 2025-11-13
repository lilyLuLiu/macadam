package e2e

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)


var _ = Describe("provider hyperv test", Label("hyperv", "windows"), func() {
	BeforeEach(func() {
		session := macadamTest.Macadam([]string{"list","--provider", "hyperv", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		err := json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(0))
	})

	AfterEach(func() {
		report := CurrentSpecReport()
		labels := report.Labels()
		stopCmd := "stop --provider hyperv"
		rmCmd := "rm --provider hyperv -f"
		for _, l := range labels {
			if l == "name" {
				stopCmd = "stop --provider hyperv myVM"
				rmCmd = "rm --provider hyperv myVM -f"
				break
			}
		}

		// stop the CentOS VM
		session := macadamTest.Macadam(strings.Fields(stopCmd))
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(session.OutputToString()).Should(ContainSubstring("stopped successfully"))

		// rm the CentOS VM and verify that "list" does not return any vm
		session = macadamTest.Macadam(strings.Fields(rmCmd))
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))

		session = macadamTest.Macadam([]string{"list", "--provider", "hyperv", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(0))
	})

	It("init CentOS VM with cpu, disk and memory setup", Label("cpu"), func() {
		// init a CentOS VM with cpu and disk-size setup
		session := macadamTest.Macadam([]string{"init", "--provider", "hyperv", "--cpus", "3", "--disk-size", "30", "--memory", "2048", image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))

		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list", "--provider", "hyperv", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM
		session = macadamTest.Macadam([]string{"start", "--provider", "hyperv"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "nproc"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(strings.TrimSpace(session.OutputToString())).Should(Equal("3"))

		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "lsblk"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(session.OutputToString()).Should(ContainSubstring("30G"))

		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "free", "-h"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		// Verify memory is close to 2048MB (allow for system overhead)
		output := session.OutputToString()
		Expect(output).Should(Or(ContainSubstring("1.7G"), ContainSubstring("1.8G"), ContainSubstring("1.9G"), ContainSubstring("2.0G")))
	})

	It("init CentOS VM with username and sshkey setup", Label("sshkey"), func() {
		// init a CentOS VM with cpu and disk-size setup
		session := macadamTest.Macadam([]string{"init", "--provider", "hyperv", "--username", "test", "--ssh-identity-path", keypath, image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))

		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list", "--provider", "hyperv", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM
		session = macadamTest.Macadam([]string{"start", "--provider", "hyperv"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "--username", "test", "whoami"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(strings.TrimSpace(session.OutputToString())).Should(Equal("test"))
	})

	It("init CentOS VM with cloud-init setup", Label("cloudinit"), func() {
		// init a CentOS VM with cpu and disk-size setup
		session := macadamTest.Macadam([]string{"init", "--provider", "hyperv", "--cloud-init", cloudinitPath, "--username", "macadamtest", "--ssh-identity-path", keypath, image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))

		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list", "--provider", "hyperv", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM
		session = macadamTest.Macadam([]string{"start", "--provider", "hyperv"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "whoami"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(strings.TrimSpace(session.OutputToString())).Should(Equal("macadamtest"))
		

		// wait until cloud-init finish
		Eventually(func() string {
			session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "systemctl", "status", "cloud-final"})
			session.WaitWithDefaultTimeout()

			if session.ExitCode() != 0 {
				return ""
			}
			return session.OutputToString()
		}, 20*time.Minute, 30*time.Second).Should(ContainSubstring("Active: active (exited)"))

		fmt.Println("cloud-init has finished")

		// ssh into the VM and check installed app
		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "git", "--version"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(session.OutputToString()).Should(ContainSubstring("git version"))

		// ssh into the VM and check file created
		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "ls", "/home/macadamtest/hello.txt"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
	})

	It("init CentOS VM with name", Label("name"), func() {
		// init a CentOS VM with name setup
		session := macadamTest.Macadam([]string{"init", "--provider", "hyperv", "--name", "myVM", image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))

		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list", "--provider", "hyperv", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM with set name
		session = macadamTest.Macadam([]string{"start", "--provider", "hyperv", "myVM"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh", "--provider", "hyperv", "myVM", "whoami"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))
		Expect(strings.TrimSpace(session.OutputToString())).Should(Equal("core"))
	})

})