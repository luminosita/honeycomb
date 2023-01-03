package handlers

import (
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/honeycomb/pkg/validators"
	"github.com/luminosita/honeycomb/pkg/validators/adapters"
)

type HandlerTemplate struct {
	validator validators.Validator
	Handler
}

func NewHandlerTemplate(h Handler) *HandlerTemplate {
	return &HandlerTemplate{
		validator: adapters.NewValidatorAdapter(),
		Handler:   h,
	}
}

func (bh *HandlerTemplate) Process(req *http.HttpRequest) (*http.HttpResponse, error) {
	//m := bh.Model(req)
	//
	//if m != nil {
	//	err := bh.validator.Validate(m)
	//	if err != nil {
	//		return http.BadValidation(err), errors.New("Validation Failed")
	//	}
	//}

	res, err := bh.Handle(req)
	if err != nil {
		return nil, err
	}

	return http.Ok(res.Body), nil
}
