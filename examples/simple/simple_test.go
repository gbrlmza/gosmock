package simple

import (
	"gotest.tools/assert"
	"testing"
)

func TestSimple_AddItem(t *testing.T) {
	simple := &Simple{}

	simple.Mock(simple.AddItem).
		WithResponse("123456", nil).
		Times(1)

	id, err := simple.AddItem("Item Name", 10)

	assert.Equal(t, "123456", id)
	assert.NilError(t, err)

	assert.Equal(t, true, simple.AllMocksUsed())
}
