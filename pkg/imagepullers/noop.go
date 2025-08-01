package imagepullers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()

	switch vmType {
	case define.AppleHvVirt, define.LibKrun:
		return copyFileMac(srcF, destF)
	default:
		return copyFile(srcF, destF)
	}
}

func copyFileMac(src, dest *os.File) error {
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

func convertToRaw(srcImg image.Image, dest *os.File) error {
	if err := srcImg.Readable(); err != nil {
		return fmt.Errorf("source image is not readable: %w", err)
	}

	return convert.Convert(dest, srcImg, convert.Options{})
}

func copyFile(src, dst *os.File) error {
	bufferedWriter := bufio.NewWriter(dst)
	defer bufferedWriter.Flush()

	_, err := io.Copy(bufferedWriter, src)
	return err
}
