package server

import (
	"context"
	"github.com/luminosita/honeycomb/pkg/http"
)

type ServerConfigurer interface {
	ServerConfig() *Config
}

type Config struct {
	BaseUrl string `mapstructure:"baseUrl" validate:"required"`

	LogCfg LoggerConfig `mapstructure:"logger"`
}

// Logger holds configuration required to customize logging for dex.
type LoggerConfig struct {
	// Level sets logging level severity.
	Level string `mapstructure:"level" validate:"omitempty,oneof=error debug info"`

	// Format specifies the format to be used for logging.
	Format string `mapstructure:"format" validate:"omitempty,oneof=text json"`
}

type ServerHandler interface {
	Config() ServerConfigurer
	Routes(context.Context) []*http.Route
	OverrideConfigItems() map[string]string
}
