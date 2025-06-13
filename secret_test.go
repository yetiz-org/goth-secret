package secret

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestWithDefaultSecret struct {
	DefaultSecret
	Field string
}

type TestWithoutDefaultSecret struct {
	Field string
}

type TestWithNestedDefaultSecret struct {
	Nested struct {
		DefaultSecret
	}
	Field string
}

func TestFindDefaultSecret(t *testing.T) {
	withSecret := TestWithDefaultSecret{}
	found, value := findDefaultSecret(reflect.ValueOf(withSecret))
	assert.True(t, found)
	assert.True(t, value.IsValid())

	withoutSecret := TestWithoutDefaultSecret{}
	found, value = findDefaultSecret(reflect.ValueOf(withoutSecret))
	assert.False(t, found)
	assert.False(t, value.IsValid())

	nestedSecret := TestWithNestedDefaultSecret{}
	found, value = findDefaultSecret(reflect.ValueOf(nestedSecret))
	assert.True(t, found)
	assert.True(t, value.IsValid())

	var nonStruct string
	found, value = findDefaultSecret(reflect.ValueOf(nonStruct))
	assert.False(t, found)
	assert.False(t, value.IsValid())

	ptrWithSecret := &TestWithDefaultSecret{}
	found, value = findDefaultSecret(reflect.ValueOf(ptrWithSecret).Elem())
	assert.True(t, found)
	assert.True(t, value.IsValid())
}
