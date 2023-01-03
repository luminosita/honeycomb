package validators

type Validator interface {
	Validate(any) []error
}
