package cmd

import (
	"github.com/luminosita/honeycomb/pkg/server"
	"github.com/luminosita/honeycomb/pkg/server/fiberadapter"
	"github.com/spf13/cobra"
)

func CommandServe(h server.ServerHandler) *cobra.Command {
	options := fiberadapter.Options{}

	cmd := &cobra.Command{
		Use:     "serve [flags] environment config-file-path",
		Short:   "Launch Bee",
		Example: "bee serve configs/boot.yaml",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			options.ConfigUrl = args[0]

			return fiberadapter.RunServe(&options, cmd.Flags(), h)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.BaseUrl, "baseUrl", "", "Base URL")

	return cmd
}
