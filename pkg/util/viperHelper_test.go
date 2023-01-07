package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Laza struct {
	A string `mapstructure:"a" validate:"required"`
}

func TestGood(t *testing.T) {
	laza := &Laza{
		A: "TestA",
	}

	items := map[string]string{"laza.a": "a"}

	err := OverrideConfig(func(key string) string {
		if key == "laza.a" {
			return "newA"
		}

		return ""
	}, items, laza)

	assert.Nil(t, err)
	assert.Equal(t, "newA", laza.A)
}

func TestBad(t *testing.T) {
	items := map[string]string{"laza.a": "a"}

	err := OverrideConfig(func(key string) string {
		return ""
	}, items, nil)

	assert.NotNil(t, err)
}

func TestBad2(t *testing.T) {
	err := OverrideConfig(nil, nil, nil)

	assert.NotNil(t, err)
}
