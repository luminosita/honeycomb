package http

type HttpResponse struct {
	StatusCode int
	Body       any

	Errors []error
}
