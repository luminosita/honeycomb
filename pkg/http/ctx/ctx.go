//go:generate mockgen -destination=./mocks/mock_fiberCtx.go -package=mocks . GetDocumenter
package ctx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/luminosita/honeycomb/pkg/validators"
	"github.com/luminosita/honeycomb/pkg/validators/adapters"
	"io"
	"mime/multipart"
)

type Ctx struct {
	fCtx *fiber.Ctx

	Body    []byte
	Params  map[string]string
	Headers map[string]string
	UserId  string

	validator validators.Validator
}

type JsonResponse struct {
	Body string `json:"body"`
}

type JsonError struct {
	Msg string `json:"error"`
}

func NewCtx(ctx *fiber.Ctx) *Ctx {
	return &Ctx{
		fCtx:      ctx,
		Body:      ctx.Body(),
		Params:    ctx.AllParams(),
		Headers:   ctx.GetReqHeaders(),
		validator: adapters.NewValidatorAdapter(),
	}
}

func (ctx *Ctx) Bind(obj any) error {
	err := ctx.fCtx.BodyParser(obj)
	if err != nil {
		return err
	}

	err = ctx.validator.Validate(obj)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Ctx) SendString(body string) error {
	return ctx.fCtx.Status(fiber.StatusOK).JSON(&JsonResponse{
		Body: body,
	})
}

func (ctx *Ctx) SendResponse(obj any) error {
	return ctx.fCtx.Status(fiber.StatusOK).JSON(obj)
}

func (ctx *Ctx) FormFile(key string) (*multipart.FileHeader, error) {
	return ctx.fCtx.FormFile(key)
}

func (ctx *Ctx) SendStream(filename string, reader io.Reader, size ...int) error {
	ctx.fCtx.Attachment(filename)

	if len(size) > 0 && size[0] >= 0 {
		return ctx.fCtx.SendStream(reader, size[0])
	} else {
		return ctx.fCtx.SendStream(reader)
	}

}
