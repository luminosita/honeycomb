package validators

import (
	"fmt"
)

type ValidationError struct {
	Errors []error `json:"error"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation errors: %+v", e.Errors)
}
