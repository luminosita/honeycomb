package adapters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type a struct {
	A string `validate:"required"`
	B int    `validate:"gte=0,lte=130"`
}

func TestValidateGood(t *testing.T) {
	a := &a{
		A: "v1",
		B: 123,
	}

	err := NewValidatorAdapter().Validate(a)

	assert.Nil(t, err)
}

func TestValidateBad(t *testing.T) {
	a := &a{
		A: "v1",
		B: 153,
	}

	err := NewValidatorAdapter().Validate(a)

	assert.NotNil(t, err)
}

func TestValidateBad2(t *testing.T) {
	a := &a{}

	err := NewValidatorAdapter().Validate(a)

	assert.NotNil(t, err)
}

func TestValidateBad3(t *testing.T) {
	err := NewValidatorAdapter().Validate(nil)

	assert.NotNil(t, err)
}
