//go:build darwin

package imagepullers

import (
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/machine/define"

	"github.com/lima-vm/go-qcow2reader"
	"github.com/lima-vm/go-qcow2reader/convert"
	"github.com/lima-vm/go-qcow2reader/image"
	"github.com/lima-vm/go-qcow2reader/image/qcow2"
	"github.com/lima-vm/go-qcow2reader/image/raw"
)

func imageExtension(vmType define.VMType, _ string) string {
	return "." + vmType.ImageFormat().Kind()
}

func doCopyFile(src *os.File, dest string) error {
	srcImg, err := qcow2reader.Open(src)
	if err != nil {
		return err
	}
	defer srcImg.Close()

	switch srcImg.Type() {
	case raw.Type:
		return copyFile(src, dest)
	case qcow2.Type:
		return convertToRaw(srcImg, dest)
	default:
		return fmt.Errorf("%s format not supported for conversion to raw", srcImg.Type())
	}
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
