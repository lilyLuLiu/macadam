package imagepullers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/containers/podman/v5/pkg/machine/define"

	"github.com/containers/podman/v5/pkg/machine/env"

	"github.com/lima-vm/go-qcow2reader"
	"github.com/lima-vm/go-qcow2reader/convert"
	"github.com/lima-vm/go-qcow2reader/image"
	"github.com/lima-vm/go-qcow2reader/image/qcow2"
	"github.com/lima-vm/go-qcow2reader/image/raw"
)

type NoopImagePuller struct {
	localPath   *define.VMFile
	sourceURI   string
	vmType      define.VMType
	machineName string
}

func NewNoopImagePuller(machineName string, vmType define.VMType) *NoopImagePuller {
	return &NoopImagePuller{
		machineName: machineName,
		vmType:      vmType,
	}
}

func (puller *NoopImagePuller) SetSourceURI(sourcePath string) {
	puller.sourceURI = sourcePath
}

func imageExtension(vmType define.VMType, sourceURI string) string {
	switch vmType {
	case define.WSLVirt:
		ext := filepath.Ext(sourceURI)
		if ext == ".wsl" {
			return ".wsl"
		}
		return ".tar.gz"
	case define.QemuVirt, define.HyperVVirt:
		return filepath.Ext(sourceURI)
	default:
		return "." + vmType.ImageFormat().Kind()
	}
}

func (puller *NoopImagePuller) LocalPath() (*define.VMFile, error) {
	// if localPath has already been calculated returns it
	if puller.localPath != nil {
		return puller.localPath, nil
	}

	// calculate and set localPath
	dirs, err := env.GetMachineDirs(puller.vmType)
	if err != nil {
		return nil, err
	}

	vmFile, err := dirs.DataDir.AppendToNewVMFile(fmt.Sprintf("%s-%s%s", puller.machineName, puller.vmType.String(), imageExtension(puller.vmType, puller.sourceURI)), nil)
	if err != nil {
		return nil, err
	}
	puller.localPath = vmFile
	return vmFile, nil
}

/*
The noopImageBuilder does not actually download any image as the image is already stored locally.
The download func is used to make a copy of the source image so that the user image is not modified
by macadam
*/
func (puller *NoopImagePuller) Download() error {
	localPath, err := puller.LocalPath()
	if err != nil {
		return err
	}
	return doCopyFile(puller.sourceURI, localPath.Path, puller.vmType)
}

func doCopyFile(src, dest string, vmType define.VMType) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	switch vmType {
	case define.AppleHvVirt, define.LibKrun:
		return copyFileMac(srcF, dest)
	case define.WSLVirt, define.HyperVVirt:
		return copyFileWin(srcF, dest)
	default:
		return copyFile(srcF, dest)
	}
}

func copyFileMac(src *os.File, dest string) error {
	srcImg, err := qcow2reader.Open(src)
	if err != nil {
		return err
	}
	defer srcImg.Close()

	switch srcImg.Type() {
	case raw.Type:
		// if the image is raw it performs a simple copy
		return copyFile(src, dest)
	case qcow2.Type:
		// if the image is qcow2 it performs a conversion to raw
		return convertToRaw(srcImg, dest)
	default:
		return fmt.Errorf("%s format not supported for conversion to raw", srcImg.Type())
	}
}

func copyFileWin(srcF *os.File, dest string) error {
	binary, err := exec.LookPath("robocopy")
	if err != nil {
		logrus.Debugf("warning: failed to get robocopy binary: %v. Falling back to default file copy between %s and %s\n", err, srcF.Name(), dest)
		return copyFile(srcF, dest)
	}

	srcDir, srcFile := filepath.Split(srcF.Name())
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

func convertToRaw(srcImg image.Image, dest string) error {
	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()

	if err := srcImg.Readable(); err != nil {
		return fmt.Errorf("source image is not readable: %w", err)
	}

	return convert.Convert(destF, srcImg, convert.Options{})
}

func copyFile(src *os.File, dest string) error {
	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()

	bufferedWriter := bufio.NewWriter(destF)
	defer bufferedWriter.Flush()

	_, err = io.Copy(bufferedWriter, src)
	return err
}
