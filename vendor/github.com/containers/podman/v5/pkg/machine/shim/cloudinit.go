package shim

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/podman/v5/pkg/machine/define"
	"github.com/containers/podman/v5/pkg/machine/vmconfigs"
)

func processOneFile(config *vmconfigs.CloudInitConfig, kind string, file string) error {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cloud-init: file %s not found: %w", file, err)
	}
	if err != nil {
		return fmt.Errorf("cloud-init: failed to access %s: %w", file, err)
	}

	switch kind {
	case "user-data":
		config.UserData = &define.VMFile{
			Path: file,
		}
	case "meta-data":
		config.MetaData = &define.VMFile{
			Path: file,
		}
	case "network-config":
		config.NetworkConfig = &define.VMFile{
			Path: file,
		}
	default:
		return fmt.Errorf("cloud-init: unexpected configuration file '%s'", kind)
	}
	return nil
}
func CmdLineCloudInitToConfig(params []string) (vmconfigs.CloudInitConfig, error) {
	config := vmconfigs.CloudInitConfig{
		UserData:      nil,
		MetaData:      nil,
		NetworkConfig: nil,
	}

	for _, param := range params {
		var kind, file string
		subparams := strings.SplitN(param, "=", 2)
		if len(subparams) == 2 {
			kind = subparams[0]
			file = subparams[1]
		} else { // len(subparams) == 1
			kind = filepath.Base(param)
			file = param
		}
		if err := processOneFile(&config, kind, file); err != nil {
			return config, err
		}
	}

	return config, nil
}
