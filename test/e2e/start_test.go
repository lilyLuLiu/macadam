package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Macadam starts", func() {

	It("non-existing VM", func() {
		session := macadamTest.Macadam([]string{"start", "123"})
		session.WaitWithDefaultTimeout()
		Expect(session).Should(gexec.Exit())
		Expect(session.ErrorToString()).Should(Equal("VM 123 does not exist"))
	})

})
