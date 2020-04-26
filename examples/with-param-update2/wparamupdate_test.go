package wparam

import (
	"gotest.tools/assert"
	"testing"
)

func TestParam_AddItem(t *testing.T) {
	simple := &MyStruct{}

	str1 := "str1"
	simple.Mock(simple.AddItem).
		WithParams("Name1", 10, &str1).
		WithParamUpdate(2, "newValue").
		WithResponse(nil).
		Times(1)

	myStr := "str1"
	err := simple.AddItem("Name1", 10, &myStr)
	assert.NilError(t, err)
	assert.Equal(t, "newValue", myStr)

	// An extra check to ensure that all mocks were used
	assert.Equal(t, true, simple.AllMocksUsed())
}

func TestParam_AddItem2(t *testing.T) {
	simple := &MyStruct{}

	simple.Mock(simple.AddItem).
		WithParams("Name1", 10, nil).
		WithResponse(nil).
		Times(1)

	err := simple.AddItem("Name1", 10, nil)
	assert.NilError(t, err)

	// An extra check to ensure that all mocks were used
	assert.Equal(t, true, simple.AllMocksUsed())
}
