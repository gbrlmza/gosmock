package wparam

import (
	"gotest.tools/assert"
	"testing"
)

func TestParam_AddItem(t *testing.T) {
	simple := &Param{}

	simple.Mock(simple.AddItem).
		WithParams("Name1", 10).
		WithResponse("100", nil).
		Times(1)

	simple.Mock(simple.AddItem).
		WithParams("Name2", 20).
		WithResponse("200", nil).
		Times(1)

	// Since we are using mocks with params, the mocked response for a given set of param
	// takes precedence over the order ot the mocked response.
	id, err := simple.AddItem("Name2", 20)
	assert.Equal(t, "200", id)
	assert.NilError(t, err)

	id, err = simple.AddItem("Name1", 10)
	assert.Equal(t, "100", id)
	assert.NilError(t, err)

	// An extra check to ensure that all mocks were used
	assert.Equal(t, true, simple.AllMocksUsed())
}
