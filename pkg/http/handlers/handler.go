package handlers

import "github.com/luminosita/honeycomb/pkg/http"

type Handler interface {
	Handle(req *http.HttpRequest) (*http.HttpResponse, error)
}
