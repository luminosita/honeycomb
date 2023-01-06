package http

import (
	"github.com/luminosita/honeycomb/pkg/http/handlers"
)

type RouteType int

const (
	STATIC RouteType = iota
	GET
	POST
	PUT
	PATCH
)

func (m RouteType) String() string {
	return []string{"STATIC", "GET", "PUT", "HEAD", "PATCH"}[m]
}

type Route struct {
	Type    RouteType
	Path    string
	Handler handlers.Handler
}
