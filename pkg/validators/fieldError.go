package validators

import (
	"fmt"
)

type FieldError struct {
	FailedField string `json:"failedField"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("Error: %s, %s, %s", e.FailedField, e.Tag, e.Value)
}
