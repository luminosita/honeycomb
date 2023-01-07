package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/luminosita/honeycomb/pkg/http"
	ctx2 "github.com/luminosita/honeycomb/pkg/http/ctx"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
	"github.com/luminosita/honeycomb/pkg/validators"
	"net/url"
)

func SetupRoute(app *fiber.App, baseUrl string, r *http.Route) error {
	path, err := url.JoinPath(baseUrl, r.Path)

	if err != nil {
		return err
	}

	switch r.Type {
	case http.STATIC:
		app.Static(r.Path, "web/static")
	case http.GET:
		app.Get(path, convert(r.Handler))
	case http.POST:
		app.Post(path, convert(r.Handler))
	case http.PUT:
		app.Put(path, convert(r.Handler))
	case http.PATCH:
		app.Patch(path, convert(r.Handler))
	}

	return nil
}

func convert(handler handlers.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newCtx := ctx2.NewCtx(ctx)

		err := handler.Handle(newCtx)
		if err != nil {
			var e *validators.ValidationError
			if errors.As(err, &e) {
				return ctx.Status(fiber.StatusBadRequest).JSON(e)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(&ctx2.JsonError{err.Error()})
		}

		return nil
	}
}
