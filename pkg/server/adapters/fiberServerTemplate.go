package adapters

import (
	"context"
	"errors"
	"fmt"
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/honeycomb/pkg/http/ctx"
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

	//	setupMiddlewares(app);

	err = bs.setupRoutes(ctx)
	if err != nil {
		return err
	}

	// This is required!!!
	bs.RefreshFiberRoutes()

	return nil
}

func (bs *FiberServerTemplate) setupRoutes(c context.Context) error {
	routes := bs.handler.Routes(c)

	for _, v := range routes {
		path, err := url.JoinPath(bs.c.BaseUrl, v.Path)

		if err != nil {
			return err
		}

		switch v.Method {
		case http.GET:
			bs.App.Get(path, ctx.Convert(v.Handler))
		case http.POST:
			bs.App.Post(path, ctx.Convert(v.Handler))
		case http.PUT:
			bs.App.Put(path, ctx.Convert(v.Handler))
		case http.PATCH:
			bs.App.Patch(path, ctx.Convert(v.Handler))
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
