package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type NoDefaultSecret struct {
	SomeField string
}

func (n *NoDefaultSecret) Name() string {
	return ""
}

func (n *NoDefaultSecret) Path() string {
	return ""
}

func TestLoadErrorCases(t *testing.T) {
	var s NoDefaultSecret
	err := Load("test", "name", &s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "struct should have a DefaultSecret field")

	oldPath := PATH
	PATH = "/non-existent-path/"
	defer func() { PATH = oldPath }()

	db := &Database{}
	err = Load("test", "name", db)
	assert.Error(t, err)
}
