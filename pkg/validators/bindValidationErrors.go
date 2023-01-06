package validators

import (
	"fmt"
)

type BindValidationErrors struct {
	Errors []error
}

func (e *BindValidationErrors) Error() string {
	return fmt.Sprintf("Validation Errors")
}
