package validators

import (
	"fmt"
)

type ValidationError struct {
	Errors []error `json:"error"`
}

func (e *ValidationError) Error() string {
	//TODO: Externalize
	return fmt.Sprintf("Validation Errors: %+v", e.Errors)
}
