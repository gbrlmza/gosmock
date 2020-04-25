package wparam

import (
	"errors"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func TestParam_AddItem(t *testing.T) {
	simple := &MyStruct{}

	simple.Mock(simple.AddItem).
		WithParams("Name1", 10).
		WithResponse(&Response{
			Err: errors.New("some_error"),
		}).Times(1)

	expectedResp := Response{Err: errors.New("some_error")}
	resp := simple.AddItem("Name1", 10)
	assert.Equal(t, true, reflect.DeepEqual(expectedResp, *resp))

	// An extra check to ensure that all mocks were used
	assert.Equal(t, true, simple.AllMocksUsed())
}
