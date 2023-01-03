package adapters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
)

type errorResponse struct {
	error string
}

func Convert(handler handlers.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		response, err := handlers.NewHandlerTemplate(handler).Process(
			&http.HttpRequest{
				Body:    ctx.Body(),
				Params:  ctx.AllParams(),
				Headers: ctx.GetReqHeaders(),
			})

		if err != nil {
			return err
		}

		ctx.SendStatus(response.StatusCode)
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
			return ctx.JSON(response.Body)
		} else {
			return ctx.JSON(response.Errors)
		}
	}
}
