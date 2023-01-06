package cmd

import (
	"context"
	"errors"
	"github.com/luminosita/honeycomb/pkg/server"
	"github.com/luminosita/honeycomb/pkg/server/adapters"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const CONFIG_ENTRY = "config"

type ServerOptions struct {
	// Flags
	BaseUrl   string
	ConfigUrl string
}

func CommandServe(h server.ServerHandler) *cobra.Command {
	options := ServerOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] environment config-file-path",
		Short:   "Launch Bee",
		Example: "bee serve configs/boot.yaml",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			options.ConfigUrl = args[0]

			return runServe(&options, cmd.Flags(), h)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.BaseUrl, "baseUrl", "", "Base URL")

	return cmd
}

func runServe(options *ServerOptions, pflags *pflag.FlagSet, h server.ServerHandler) error {
	ctx := context.Background()

	bootData, err := os.ReadFile(options.ConfigUrl)
	if err != nil {
		return err
	}

	boot := rkboot.NewBoot(rkboot.WithBootConfigRaw(bootData))
	boot.Bootstrap(ctx)

	vpr, err := setupViper(options, pflags)
	if err != nil {
		return err
	}

	srv := adapters.NewFiberServerTemplate(h)

	err = srv.Run(ctx, vpr)
	if err != nil {
		return err
	}

	boot.WaitForShutdownSig(ctx)

	return nil
}

func setupViper(options *ServerOptions, pflags *pflag.FlagSet) (*viper.Viper, error) {
	vpr := rkentry.GlobalAppCtx.GetConfigEntry(CONFIG_ENTRY)

	if vpr == nil {
		//TODO: Externalize
		return nil, errors.New("Unable to load connfiguration. Check the configuration file path")
	}

	if options.BaseUrl != "" {
		err := viper.BindPFlag("config.server.baseUrl", pflags.Lookup("baseUrl"))
		if err != nil {
			return nil, err
		}
	}

	//TOOD: Environnment vars are not working
	viper.SetEnvPrefix("bee") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return vpr.Viper, nil
}
