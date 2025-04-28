package main

import (
	"github.com/crc-org/macadam/cmd/macadam/registry"
	macadam "github.com/crc-org/macadam/pkg/machinedriver"
	"github.com/spf13/cobra"
)

var (
	stopCmd = &cobra.Command{
		Use:     "stop [options] [MACHINE]",
		Short:   "Stop an existing machine",
		Long:    "Stop a managed virtual machine ",
		RunE:    stop,
		Args:    cobra.MaximumNArgs(1),
		Example: `macadam stop`,
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: stopCmd,
	})
}

func stop(cmd *cobra.Command, args []string) error {
	machineName := defaultMachineName
	if len(args) > 0 && len(args[0]) > 0 {
		machineName = args[0]
	}

	driver, err := macadam.GetDriverByMachineName(machineName)
	if err != nil {
		return err
	}

	return driver.Stop()
}
