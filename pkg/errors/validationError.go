package errors

import (
	"fmt"
)

type ValidationError struct {
	FailedField string `json:"failedField"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Error: %s, %s, %s", e.FailedField, e.Tag, e.Value)
}
