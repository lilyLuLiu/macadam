package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	macadamTest *MacadamTestIntegration

	_ = BeforeEach(func() {
		macadamTest = MacadamTestCreate()
	})
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Macadam suite")
}
