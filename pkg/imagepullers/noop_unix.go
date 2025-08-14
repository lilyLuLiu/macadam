//go:build !windows && !darwin

package imagepullers

import (
	"os"
	"path/filepath"

	"github.com/containers/podman/v5/pkg/machine/define"
)

func imageExtension(_ define.VMType, sourceURI string) string {
	return filepath.Ext(sourceURI)
}

func doCopyFile(src *os.File, dest string) error {
	return copyFile(src, dest)
}
