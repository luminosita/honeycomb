package http

import (
	"github.com/luminosita/honeycomb/pkg/http/handlers"
)

type Method int

const (
	GET Method = iota
	POST
	PUT
	PATCH
)

func (m Method) String() string {
	return []string{"GET", "PUT", "HEAD", "PATCH"}[m]
}

type Route struct {
	Method  Method
	Path    string
	Handler handlers.Handler
}
