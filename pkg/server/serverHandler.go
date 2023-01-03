package server

import (
	"context"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
)

type Method int
type Environment int

const (
	GET   Method = iota // Head = 0
	POST                // Shoulder = 1
	PUT                 // Knee = 2
	PATCH               // Toe = 3
)

const (
	DEV   Environment = iota // Head = 0
	STAGE                    // Shoulder = 1
	PROD                     // Knee = 2
)

func (m Method) String() string {
	return []string{"GET", "PUT", "HEAD", "PATCH"}[m]
}
func (e Environment) String() string {
	return []string{"dev", "stage", "prod"}[e]
}

func EnvironmentFromString(str string) Environment {
	return map[string]Environment{"dev": DEV, "stage": STAGE, "prod": PROD}[str]
}

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

type Route struct {
	Method  Method
	Path    string
	Handler handlers.Handler
}

type ServerHandler interface {
	Config() ServerConfigurer
	Routes(context.Context) []*Route
	OverrideConfigItems() map[string]string
}
