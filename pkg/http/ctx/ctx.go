package ctx

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
	"github.com/luminosita/honeycomb/pkg/validators"
	"github.com/luminosita/honeycomb/pkg/validators/adapters"
)

type Ctx struct {
	*fiber.Ctx

	Body    []byte
	Params  map[string]string
	Headers map[string]string
	UserId  string

	validator validators.Validator
}

type JsonError struct {
	error error
}

func Convert(handler handlers.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newCtx := NewCtx(ctx)

		err := handler(newCtx)
		if err != nil {
			var e *BindValidationErrors
			if errors.As(err, &e) {
				return ctx.Status(fiber.StatusBadRequest).JSON(e.errs)
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(&JsonError{err})
		}

		return nil
	}
}

func NewCtx(ctx *fiber.Ctx) *Ctx {
	return &Ctx{
		Ctx:       ctx,
		Body:      ctx.Body(),
		Params:    ctx.AllParams(),
		Headers:   ctx.GetReqHeaders(),
		validator: adapters.NewValidatorAdapter(),
	}
}

func (ctx *Ctx) Bind(obj any) error {
	err := ctx.BodyParser(obj)
	if err != nil {
		return err
	}

	v_errs := ctx.validator.Validate(obj)
	if v_errs != nil {
		return &BindValidationErrors{errs: v_errs}
	}

	return nil
}

func (ctx *Ctx) SendResponse(obj any) error {
	return ctx.Status(fiber.StatusOK).JSON(obj)
}
