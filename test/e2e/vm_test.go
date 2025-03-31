package e2e

import (
	"encoding/json"
	"os"

	"github.com/crc-org/macadam/test/osprovider"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type ListReporter struct {
	Image          string
	Created        string
	Running        bool
	Starting       bool
	LastUp         string
	CPUs           uint64
	Memory         string
	DiskSize       string
	Port           int
	RemoteUsername string
	IdentityPath   string
	VMType         string
}

var _ = Describe("Macadam", func() {

	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "test-")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	It("creates a new CentOS VM, starts it, ssh in and cleans", func() {
		// verify there is no vm
		var machineResponses []ListReporter
		session := macadamTest.Macadam([]string{"list"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		list := session.OutputToString()
		Expect(list).Should(Equal(""))

		// download CentOS image
		centosProvider := osprovider.NewCentosProvider()
		image, err := centosProvider.Fetch(tempDir)
		Expect(err).NotTo(HaveOccurred())

		// init a CentOS VM
		session = macadamTest.Macadam([]string{"init", image})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())

		// check the list command returns one item
		session = macadamTest.Macadam([]string{"list"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		list = session.OutputToString()
		err = json.Unmarshal([]byte(list), &machineResponses)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(machineResponses)).Should(Equal(1))

		// start the CentOS VM
		session = macadamTest.Macadam([]string{"start"})
		session.WaitWithTimeout(180)
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(ContainSubstring("started successfully"))

		// ssh into the VM and prints user
		session = macadamTest.Macadam([]string{"ssh", "whoami"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(Equal("core"))

		// stop the CentOS VM
		session = macadamTest.Macadam([]string{"stop"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.OutputToString()).Should(ContainSubstring("stopped successfully"))

		// rm the CentOS VM and verify that "list" does not return any vm
		session = macadamTest.Macadam([]string{"rm", "-f"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())

		session = macadamTest.Macadam([]string{"list"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		list = session.OutputToString()
		Expect(list).Should(Equal(""))
	})

})
