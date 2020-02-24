# gosmock: go simple mock

A simple & easy to use mocking tool

## Why?

Sometimes you don't need a whole testing framework like [mock](https://github.com/golang/mock) or [testify](https://github.com/stretchr/testify). Also was used to learn & practice some concepts of [reflection](https://golang.org/pkg/reflect/) in go.

## Examples

### Basic

The mocked interface

```go
package main 

import "github.com/gbrlmza/gosmock"

// The mocked implementation
type Simple struct {
	gosmock.MockTool
}

func (s *Simple) AddItem(name string, quantity int) (p1 string, p2 error) {
	s.GetMockedResponse(s.AddItem, name, quantity).Fill(&p1, &p2)
	return
}
```

The test

```go
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

	assert.Equal(t, 0, simple.UnusedMocks())
}

```

### Specifying params

Test for the same struct of the previous example

```go
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
	assert.Equal(t, 0, simple.UnusedMocks())
}
``` 