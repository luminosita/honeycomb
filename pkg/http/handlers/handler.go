package handlers

import (
	"github.com/luminosita/honeycomb/pkg/http/ctx"
)

type Handler interface {
	Handle(ctx *ctx.Ctx) error
}
