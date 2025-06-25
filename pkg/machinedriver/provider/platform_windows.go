package provider

import (
	"fmt"

	"github.com/containers/podman/v5/pkg/machine/define"
	hypervPkg "github.com/containers/podman/v5/pkg/machine/hyperv"
	"github.com/containers/podman/v5/pkg/machine/vmconfigs"
	wslPkg "github.com/containers/podman/v5/pkg/machine/wsl"
)

const wsl = "wsl"
const hyperv = "hyperv"

func GetProviderOrDefault(name string) (vmconfigs.VMProvider, error) {
	if name == "" {
		name = define.WSLVirt.String()
	}
	switch name {
	case wsl:
		return new(wslPkg.WSLStubber), nil
	case hyperv:
		return new(hypervPkg.HyperVStubber), nil
	default:
		return nil, fmt.Errorf("unknown provider `%s`. Valid providers are: %v", name, GetProviders())
	}
}

func GetProviders() []string {
	return []string{wsl, hyperv}
}

func GetDefaultProvider() string {
	return wsl
}
