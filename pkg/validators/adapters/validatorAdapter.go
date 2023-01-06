package adapters

import (
	errors2 "errors"
	"github.com/go-playground/validator/v10"
	"github.com/luminosita/honeycomb/pkg/log"
	"github.com/luminosita/honeycomb/pkg/validators"
)

type ValidatorAdapter struct {
	validator *validator.Validate
}

func NewValidatorAdapter() *ValidatorAdapter {
	return &ValidatorAdapter{
		validator: validator.New(),
	}
}

func (v *ValidatorAdapter) Validate(obj any) []error {
	var e []error

	log.Log().Debugf("Validation input: %+v", obj)

	err := v.validator.Struct(obj)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			//TODO: Externalize
			return []error{errors2.New("Invalid validation value")}
		}

		for _, err := range err.(validator.ValidationErrors) {
			element := validators.ValidationError{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}

			e = append(e, &element)
		}
	}
	return e
}
