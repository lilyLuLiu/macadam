package utils

import (
	"github.com/onsi/gomega/types"
)

type ExitMatcher struct {
	types.GomegaMatcher
	ExpectedExitCode int
	ExitCode         int
	ExpectedStderr   string
}

// ExitWithError checks both exit code and stderr, fails if either does not match
// Modeled after the gomega Exit() matcher and also operates on sessions.
func ExitWithError(expectExitCode int, expectStderr string) *ExitMatcher {
	return &ExitMatcher{ExpectedExitCode: expectExitCode, ExpectedStderr: expectStderr}
}
