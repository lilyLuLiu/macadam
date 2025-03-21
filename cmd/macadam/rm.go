package main

import (
	"github.com/containers/podman/v5/pkg/machine"
	"github.com/crc-org/macadam/cmd/macadam/registry"
	macadam "github.com/crc-org/macadam/pkg/machinedriver"
	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:     "rm [options] [MACHINE]",
		Short:   "Remove an existing machine",
		Long:    "Remove a managed virtual machine ",
		RunE:    rm,
		Args:    cobra.MaximumNArgs(1),
		Example: `macadam rm`,
	}
)

var (
	destroyOptions machine.RemoveOptions
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: rmCmd,
	})

	flags := rmCmd.Flags()
	formatFlagName := "force"
	flags.BoolVarP(&destroyOptions.Force, formatFlagName, "f", false, "Stop and do not prompt before rming")
}

func rm(_ *cobra.Command, args []string) error {
	machineName := defaultMachineName
	if len(args) > 0 && len(args[0]) > 0 {
		machineName = args[0]
	}

	driver, err := macadam.GetDriverByMachineName(machineName)
	if err != nil {
		return nil
	}

	return driver.RemoveWithOptions(destroyOptions)
}
