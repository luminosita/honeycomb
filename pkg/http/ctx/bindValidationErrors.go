package ctx

import (
	"fmt"
)

type BindValidationErrors struct {
	errs []error
}

func (e *BindValidationErrors) Error() string {
	return fmt.Sprintf("Validation Errors")
}
