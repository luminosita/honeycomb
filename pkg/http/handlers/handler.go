package handlers

import (
	"github.com/luminosita/honeycomb/pkg/http/ctx"
)

type Handler = func(ctx *ctx.Ctx) error
