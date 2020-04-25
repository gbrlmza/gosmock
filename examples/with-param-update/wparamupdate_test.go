package wparam

import (
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func TestParam_AddItem(t *testing.T) {
	simple := &MyStruct{}

	p := []int{0}
	simple.Mock(simple.AddItem).
		WithParams("Name1", 10, &p).
		WithParamUpdate(2, []int{1, 2, 3, 4}).
		WithResponse(nil).
		Times(1)

	err := simple.AddItem("Name1", 10, &p)
	assert.NilError(t, err)
	assert.Equal(t, true, reflect.DeepEqual(p, []int{1, 2, 3, 4}))

	// An extra check to ensure that all mocks were used
	assert.Equal(t, true, simple.AllMocksUsed())
}
