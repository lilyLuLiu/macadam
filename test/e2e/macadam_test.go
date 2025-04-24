package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega/gexec"
)

// MacadamTestIntegration struct for command line options
type MacadamTestIntegration struct {
	MacadamBinary string
}

type MacadamExecOptions struct {
	Wrapper []string // A command to run, receiving the Macadam command line. default: none
}

func getMacadamBinary(cwd string) string {
	macadamBinary := filepath.Join(cwd, "..", "..", "bin", fmt.Sprintf("macadam-%s-%s", runtime.GOOS, runtime.GOARCH))
	return macadamBinary
}

// MacadamTestCreate creates a MacadamTestIntegration instance for the tests
func MacadamTestCreate() *MacadamTestIntegration {
	cwd, _ := os.Getwd()

	macadamBinary := getMacadamBinary(cwd)

	return &MacadamTestIntegration{
		MacadamBinary: macadamBinary,
	}
}

// Macadam executes macadam on the filesystem with default options.
func (m *MacadamTestIntegration) Macadam(args []string) *MacadamSession {
	var command *exec.Cmd
	macadamBinary := m.MacadamBinary

	command = exec.Command(macadamBinary, args...)

	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	if err != nil {
		Fail(fmt.Sprintf("unable to run macadam command: %s\n%v", strings.Join(args, " "), err))
	}
	return &MacadamSession{session}
}
