package e2e

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"os"
	"os/exec"

	"github.com/crc-org/macadam/test/osprovider"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var tempDir string
var machineResponses []ListReporter
var err error
var image string
var keypath string
var cloudinitPath string

var _ = BeforeSuite(func() {
    tempDir, err = os.MkdirTemp("", "test-")
	Expect(err).NotTo(HaveOccurred())

	// download CentOS image
	centosProvider := osprovider.NewCentosProvider()
	image, err = centosProvider.Fetch(tempDir)
	fmt.Println(image)
	Expect(err).NotTo(HaveOccurred())
	
	keypath := filepath.Join(tempDir,"id_rsa")
	cloudinitPath := filepath.Join(tempDir,"user-data")
	//generate ssh key
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-f", keypath, "-N", "")
	err := cmd.Run()
	Expect(err).ShouldNot(HaveOccurred())
	//copy user-data
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get working directory: %v\n", err)
	}
	cloudinit := wd+"/../testdata/user-data"
	content, err := os.ReadFile(cloudinit)
	if err != nil {
		fmt.Printf("failed to read %s: %v\n", cloudinit, err)
	}
	err = os.WriteFile(cloudinitPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("failed to write to %s: %v\n", cloudinitPath, err)
	}
})

var _ = AfterSuite(func() {
    os.RemoveAll(tempDir)
})

 
var _ = Describe("Macadam init setup test", Label("init"), func() {
	BeforeEach(func() {
		session := macadamTest.Macadam([]string{"list", "--format", "json"})
		Eventually(session).Should(gexec.Exit(0))
		err := json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(0))
	})

	AfterEach(func() {
		// stop the CentOS VM
		session := macadamTest.Macadam([]string{"stop"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(ContainSubstring("stopped successfully"))

		// rm the CentOS VM and verify that "list" does not return any vm
		session = macadamTest.Macadam([]string{"rm", "-f"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())

		session = macadamTest.Macadam([]string{"list", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(0))	
	})

	It("init CentOS VM with cpu, disk and memory setup", Label("cpu"), func() {
		// init a CentOS VM with cpu and disk-size setup
		session := macadamTest.Macadam([]string{"init", "--cpus", "3", "--disk-size", "30", "--memory","2048", image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))


		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM
		session = macadamTest.Macadam([]string{"start"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh", "nproc"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(Equal("3"))

		session = macadamTest.Macadam([]string{"ssh", "lsblk"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(ContainSubstring("30G"))

		session = macadamTest.Macadam([]string{"ssh", "free", "-h"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		fmt.Println(session.OutputToString())
		Expect(session.OutputToString()).Should(ContainSubstring("1.7"))
	})

	It("init CentOS VM with username and sshkey setup", Label("test"), func() {
		// init a CentOS VM with cpu and disk-size setup
		session := macadamTest.Macadam([]string{"init", "--username", "test","--ssh-identity-path", keypath, image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))

		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM
		session = macadamTest.Macadam([]string{"start"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh","--username", "test", "whoami"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(Equal("test"))
	})

	It("init CentOS VM with cloud-init setup", Label("cloudinit"), func() {
		// init a CentOS VM with cpu and disk-size setup
		session := macadamTest.Macadam([]string{"init", "--cloud-init", cloudinitPath,"--username", "macadamtest", "--ssh-identity-path", keypath, image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit(0))


		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list", "--format", "json"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		err = json.Unmarshal(session.Out.Contents(), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM
		session = macadamTest.Macadam([]string{"start"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh", "whoami"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(Equal("macadamtest"))
	})

})
