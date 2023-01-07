package adapters

import (
	errors2 "errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/luminosita/honeycomb/pkg/validators"
	"github.com/pkg/errors"
)

type ValidatorAdapter struct {
	validator *validator.Validate
}

func NewValidatorAdapter() *ValidatorAdapter {
	return &ValidatorAdapter{
		validator: validator.New(),
	}
}

func (v *ValidatorAdapter) Validate(obj any) error {
	if obj == nil {
		return errors.New(fmt.Sprintf("Bad validation request: %+v", obj))
	}

	var e []error

	err := v.validator.Struct(obj)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			//TODO: Externalize
			return errors2.New(fmt.Sprintf("Invalid validation value: %+v", err))
		}

		for _, err := range err.(validator.ValidationErrors) {
			element := validators.FieldError{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}

			e = append(e, &element)
		}
	}

	if e != nil {
		return &validators.ValidationError{Errors: e}
	}
	
	return nil
}
