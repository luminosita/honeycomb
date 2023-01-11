package server

import (
	"context"
	"errors"
	"github.com/luminosita/honeycomb/pkg/server/adapters"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	rkgrpc "github.com/rookie-ninja/rk-grpc/v2/boot"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const ConfigEntry = "config"
const GrpcConfigEntry = "grpc"

type Options struct {
	// Flags
	BaseUrl   string
	ConfigUrl string
}

func RunServe(options *Options, pflags *pflag.FlagSet, h ServerHandler) error {
	ctx := context.Background()

	bootData, err := os.ReadFile(options.ConfigUrl)
	if err != nil {
		return err
	}

	boot := rkboot.NewBoot(rkboot.WithBootConfigRaw(bootData))

	var grpcEntry *rkgrpc.GrpcEntry

	if grpcHandler, ok := h.(GrpcHandler); ok {
		grpcEntry = rkgrpc.GetGrpcEntry(GrpcConfigEntry)

		//	rkgrpcjwt.UnaryServerInterceptor()

		grpcEntry.AddRegFuncGrpc(grpcHandler.GrpcRegFunc)
		grpcEntry.AddRegFuncGw(grpcHandler.GwRegFunc)
	}

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

	if grpcEntry != nil {
		grpcEntry.Interrupt(context.Background())
	}

	return nil
}

func setupViper(options *Options, pflags *pflag.FlagSet) (*viper.Viper, error) {
	vpr := rkentry.GlobalAppCtx.GetConfigEntry(ConfigEntry)

	if vpr == nil {
		return nil, errors.New("unable to load configuration. Check the configuration file path")
	}

	if options.BaseUrl != "" {
		err := vpr.BindPFlag("config.server.baseUrl", pflags.Lookup("baseUrl"))
		if err != nil {
			return nil, err
		}
	}

	vpr.SetEnvPrefix("bee") // will be uppercased automatically
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vpr.AutomaticEnv()

	return vpr.Viper, nil
}
