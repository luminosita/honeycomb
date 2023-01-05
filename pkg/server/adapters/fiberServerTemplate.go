package adapters

import (
	"context"
	"errors"
	"fmt"
	"github.com/luminosita/honeycomb/pkg/http/adapters"
	"github.com/luminosita/honeycomb/pkg/log"
	"github.com/luminosita/honeycomb/pkg/server"
	adapters2 "github.com/luminosita/honeycomb/pkg/validators/adapters"
	rkfiber "github.com/rookie-ninja/rk-fiber/boot"
	"github.com/spf13/viper"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

const CFG_ENTRY = "config"
const FIBER_CFG_ENTRY = "fiber"
const TAG_NAME = "mapstructure"

// Config holds the server's configuration options.
//
// Multiple servers using the same storage are expected to be configured identically.
type FiberServerTemplate struct {
	env server.Environment

	c *server.Config

	handler server.ServerHandler

	baseURL *url.URL

	*rkfiber.FiberEntry
}

// NewServer constructs a server from the provided config.
func NewFiberServerTemplate(env server.Environment, h server.ServerHandler) *FiberServerTemplate {
	return newFiberServerTemplate(env, h)
}

func newFiberServerTemplate(env server.Environment, h server.ServerHandler) *FiberServerTemplate {
	return &FiberServerTemplate{
		env:     env,
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

	//	setupMiddlewares(app);

	err = bs.setupRoutes(ctx)
	if err != nil {
		return err
	}

	// This is required!!!
	bs.RefreshFiberRoutes()

	return nil
}

func (bs *FiberServerTemplate) setupRoutes(ctx context.Context) error {
	routes := bs.handler.Routes(ctx)

	for _, v := range routes {
		path, err := url.JoinPath(bs.c.BaseUrl, v.Path)

		if err != nil {
			return err
		}

		switch v.Method {
		case server.GET:
			bs.App.Get(path, adapters.Convert(v.Handler))
		case server.POST:
			bs.App.Post(path, adapters.Convert(v.Handler))
		case server.PUT:
			bs.App.Put(path, adapters.Convert(v.Handler))
		case server.PATCH:
			bs.App.Patch(path, adapters.Convert(v.Handler))
		}
	}

	return nil
}

func (bs *FiberServerTemplate) setupLogger() log.Logger {
	log.SetLogger(bs.c.LogCfg.Level, bs.c.LogCfg.Format)

	logger := log.Log()

	logger.Infof(
		"Bee Version: %s, Go Version: %s, Go OS/ARCH: %s %s",
		bs.env,
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

	bs.overrideConfig(viper)

	return nil
}

func (bs *FiberServerTemplate) overrideConfig(viper *viper.Viper) {
	//TODO: Not working with overrides
	t := reflect.TypeOf(bs.c).Elem()
	s := reflect.ValueOf(bs.c).Elem()

	for k, v := range bs.handler.OverrideConfigItems() {
		newValue := viper.GetString(k)

		sp := strings.Split(k, ".")
		tagName := sp[len(sp)-1]

		for i := 0; i < t.NumField(); i++ {
			tv, ok := t.Field(i).Tag.Lookup(TAG_NAME)

			if ok && tv == tagName {
				f := s.FieldByName(t.Field(i).Name)
				if f.Kind() == reflect.String && f.CanSet() {
					f.SetString(newValue)
				} else {
					//TODO: Externalize
					fmt.Printf("Wrong config field to override: %s", v)
				}
			}
		}
	}
}
