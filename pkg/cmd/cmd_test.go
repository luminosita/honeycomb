package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandServer(t *testing.T) {
	assert.NotNil(t, CommandServe(nil))
}

func TestCommandVersion(t *testing.T) {
	assert.NotNil(t, CommandVersion())
}
