package provider

import (
	"fmt"

	"github.com/containers/podman/v5/pkg/machine/applehv"
	"github.com/containers/podman/v5/pkg/machine/define"
	"github.com/containers/podman/v5/pkg/machine/vmconfigs"
	"github.com/sirupsen/logrus"
)

var defaultProvider = define.AppleHvVirt

func GetProviderOrDefault(name string) (vmconfigs.VMProvider, error) {
	resolvedVMType, err := define.ParseVMType(name, defaultProvider)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Using macadam with `%s` virtualization provider", resolvedVMType.String())
	switch resolvedVMType {
	case define.AppleHvVirt:
		return new(applehv.AppleHVStubber), nil
	default:
		return nil, fmt.Errorf("unknown provider `%s`. Valid providers are: %v", name, GetProviders())
	}
}

func GetProviders() []string {
	return []string{define.AppleHvVirt.String()}
}

func GetDefaultProvider() string {
	return defaultProvider.String()
}
