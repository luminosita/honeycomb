package adapters

import (
	"context"
	"errors"
	"github.com/luminosita/common-bee/pkg/log"
	"github.com/luminosita/honeycomb/pkg/server"
	"github.com/luminosita/honeycomb/pkg/server/middleware"
	"github.com/luminosita/honeycomb/pkg/utils"
	adapters2 "github.com/luminosita/honeycomb/pkg/validators/adapters"
	rkfiber "github.com/rookie-ninja/rk-fiber/boot"
	"github.com/spf13/viper"
	"net/url"
	"runtime"
)

const CFG_ENTRY = "config"
const FIBER_CFG_ENTRY = "fiber"

const JWT_ENV_KEY = "jwt_secret"

type FiberServerTemplate struct {
	c *server.Config

	handler server.ServerHandler

	baseURL *url.URL

	*rkfiber.FiberEntry
}

func NewFiberServerTemplate(h server.ServerHandler) *FiberServerTemplate {
	return newFiberServerTemplate(h)
}

func newFiberServerTemplate(h server.ServerHandler) *FiberServerTemplate {
	return &FiberServerTemplate{
		handler: h,
	}
}

func (bs *FiberServerTemplate) Run(ctx context.Context, viper *viper.Viper) error {
	err := bs.loadConfig(viper)
	if err != nil {
		return err
	}

	bs.setupLogger()

	bs.baseURL, err = url.Parse(bs.c.BaseUrl)
	if err != nil {
		return err
	}

	bs.FiberEntry = rkfiber.GetFiberEntry(FIBER_CFG_ENTRY)

	err = bs.setupMiddlewares(viper)
	if err != nil {
		return err
	}

	err = bs.setupRoutes(ctx)
	if err != nil {
		return err
	}

	// This is required!!!
	bs.RefreshFiberRoutes()

	return nil
}

func (bs *FiberServerTemplate) setupMiddlewares(viper *viper.Viper) error {
	a := server.AccessFromString(bs.c.Access)

	if a == server.RESTRICTED {
		path, err := url.JoinPath(bs.c.BaseUrl)
		if err != nil {
			return err
		}

		secret := viper.GetString(JWT_ENV_KEY)
		if len(secret) == 0 {
			return errors.New("JWT Secret not configured")
		}

		bs.App.Use(path, middleware.Protected(secret))
	}

	return nil
}

func (bs *FiberServerTemplate) setupRoutes(c context.Context) error {
	routes := bs.handler.Routes(c)

	for _, v := range routes {
		err := utils.SetupRoute(bs.App, bs.c.BaseUrl, v)

		if err != nil {
			return err
		}
	}

	return nil
}

func (bs *FiberServerTemplate) setupLogger() log.Logger {
	log.SetLogger(bs.c.LogCfg.Level, bs.c.LogCfg.Format)

	logger := log.Log()

	//TODO: Read version from somewhere
	version := "DEV"

	logger.Infof(
		"Bee Version: %s, Go Version: %s, Go OS/ARCH: %s %s",
		version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	return logger
}

func (bs *FiberServerTemplate) loadConfig(viper *viper.Viper) error {
	c := bs.handler.Config()

	err := viper.UnmarshalKey(CFG_ENTRY, c)
	if err != nil {
		return err
	}

	validator := adapters2.NewValidatorAdapter()

	res := validator.Validate(c)
	if res != nil {
		log.Log().Errorf("%+v", res)
		//TODO: Externalize
		return errors.New("Failed to load configuration")
	}

	bs.c = c.ServerConfig()

	return utils.OverrideConfig(viper.GetString, bs.handler.OverrideConfigItems(), bs.c)
}
