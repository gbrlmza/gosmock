package wparam

import (
	"gotest.tools/assert"
	"testing"
)

func TestParam_AddItem(t *testing.T) {
	simple := &MyStruct{}

	num := 1
	str1 := "str1"
	str2 := "str2"
	simple.Mock(simple.AddItem).
		WithParams(&str1, &num, &str2).
		WithParamUpdate(0, "newstr1").
		WithParamUpdate(1, 2).
		WithParamUpdate(2, "newstr2").
		WithResponse(nil).
		Times(1)

	err := simple.AddItem(&str1, &num, &str2)
	assert.NilError(t, err)
	assert.Equal(t, "newstr1", str1)
	assert.Equal(t, 2, num)
	assert.Equal(t, "newstr2", str2)

	// An extra check to ensure that all mocks were used
	assert.Equal(t, true, simple.AllMocksUsed())
}
