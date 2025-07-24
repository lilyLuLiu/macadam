//go:build !windows && !darwin

package provider

import (
	"fmt"

	"github.com/containers/podman/v5/pkg/machine/define"
	qemuPkg "github.com/containers/podman/v5/pkg/machine/qemu"
	"github.com/containers/podman/v5/pkg/machine/vmconfigs"
	"github.com/sirupsen/logrus"
)

var defaultProvider = define.QemuVirt

func GetProviderOrDefault(name string) (vmconfigs.VMProvider, error) {
	resolvedVMType, err := define.ParseVMType(name, defaultProvider)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Using macadam with `%s` virtualization provider", resolvedVMType.String())
	switch resolvedVMType {
	case define.QemuVirt:
		return new(qemuPkg.QEMUStubber), nil
	default:
		return nil, fmt.Errorf("unknown provider `%s`. Valid providers are: %v", name, GetProviders())
	}
}

func GetProviders() []string {
	return []string{define.QemuVirt.String()}
}

func GetDefaultProvider() string {
	return defaultProvider.String()
}
