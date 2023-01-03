package http

type HttpRequest struct {
	Body    []byte
	Params  map[string]string
	Headers map[string]string
	UserId  string
}
