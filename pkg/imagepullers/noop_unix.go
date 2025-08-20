//go:build !windows && !darwin

package imagepullers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/podman/v5/pkg/machine/define"
)

func imageExtension(_ define.VMType, sourceURI string) (string, error) {
	ext := strings.ToLower(filepath.Ext(sourceURI))
	if ext != ".qcow2" && ext != ".raw" {
		return "", fmt.Errorf("unsupported image extension %s; supported formats are .qcow2 and .raw", ext)
	}
	return ext, nil
}

func doCopyFile(src *os.File, dest string) error {
	return copyFile(src, dest)
}
