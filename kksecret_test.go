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
	// 測試直接包含 DefaultSecret 的結構體
	withSecret := TestWithDefaultSecret{}
	found, value := findDefaultSecret(reflect.ValueOf(withSecret))
	assert.True(t, found)
	assert.True(t, value.IsValid())

	// 測試不包含 DefaultSecret 的結構體
	withoutSecret := TestWithoutDefaultSecret{}
	found, value = findDefaultSecret(reflect.ValueOf(withoutSecret))
	assert.False(t, found)
	assert.False(t, value.IsValid())

	// 測試嵌套包含 DefaultSecret 的結構體
	nestedSecret := TestWithNestedDefaultSecret{}
	found, value = findDefaultSecret(reflect.ValueOf(nestedSecret))
	assert.True(t, found)
	assert.True(t, value.IsValid())

	// 測試非結構體值
	var nonStruct string
	found, value = findDefaultSecret(reflect.ValueOf(nonStruct))
	assert.False(t, found)
	assert.False(t, value.IsValid())

	// 測試指針類型
	ptrWithSecret := &TestWithDefaultSecret{}
	found, value = findDefaultSecret(reflect.ValueOf(ptrWithSecret).Elem())
	assert.True(t, found)
	assert.True(t, value.IsValid())
}
