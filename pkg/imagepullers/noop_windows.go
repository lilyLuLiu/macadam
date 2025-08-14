//go:build windows

package imagepullers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/containers/podman/v5/pkg/machine/define"
	"github.com/sirupsen/logrus"
)

func imageExtension(vmType define.VMType, sourceURI string) string {
	switch vmType {
	case define.WSLVirt:
		ext := filepath.Ext(sourceURI)
		if ext == ".wsl" {
			return ".wsl"
		}
		return ".tar.gz"
	default:
		return filepath.Ext(sourceURI)
	}
}

func doCopyFile(src *os.File, dest string) error {
	binary, err := exec.LookPath("robocopy")
	if err != nil {
		logrus.Debugf("warning: failed to get robocopy binary: %v. Falling back to default file copy between %s and %s\n", err, src.Name(), dest)
		return copyFile(src, dest)
	}

	srcDir, srcFile := filepath.Split(src.Name())
	destDir := filepath.Dir(dest)

	// Executes the robocopy command with options optimized for a fast, single-file copy.
	// /J:   Copies using unbuffered I/O (better for large files).
	// /MT:  Enables multi-threaded copying for improved performance.
	// /R:0: Sets retries on failed copies to 0 to avoid long waits.
	// /IS:  Includes Same files, which forces an overwrite even if the destination
	//       file appears identical to the source.
	cmd := exec.Command(binary, "/J", "/MT", "/R:0", "/IS", srcDir, destDir, srcFile)
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	err = cmd.Run()
	if err != nil {
		// robocopy does not use classic exit codes like linux commands, so we need to check for specific errors
		// Only exit code 8 is considered an error, all other exit codes are considered successful with exceptions
		// see https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/robocopy
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode := exitErr.ExitCode()
			if exitCode >= 8 {
				return fmt.Errorf("failed to run robocopy: %w", err)
			}
		} else {
			return fmt.Errorf("failed to run robocopy: %w", err)
		}
	}

	if err := os.Remove(dest); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove existing destination file: %w", err)
	}

	err = os.Rename(filepath.Join(destDir, srcFile), dest)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}
