package gosmock

import (
	"fmt"
	"reflect"
)

type Call struct {
	mTool      *MockTool
	fn         interface{}
	fnType     reflect.Type
	fnName     string
	Count      int
	Params     []interface{}
	paramsHash string
	Response   []interface{}
}

func (c *Call) WithParams(p ...interface{}) *Call {
	c.Params = p
	c.paramsHash = hash(p)
	return c
}

func (c *Call) WithResponse(r ...interface{}) *Call {
	c.Response = r
	return c
}

func (c *Call) Times(n int) *Call {
	// Validate input params
	if c.Params != nil {
		c.compareTypes(c.fnType.NumIn(), c.fnType.In, c.Params)
	}

	// Validate output params
	c.compareTypes(c.fnType.NumOut(), c.fnType.Out, c.Response)

	for i := 0; i < n; i++ {
		c.mTool.responses[c.fnName] = append(c.mTool.responses[c.fnName], *c)
	}
	return c
}

func (c *Call) compareTypes(number int, expectedFn func(i int) reflect.Type, params []interface{}) {
	// Validate number of params
	if number != len(params) {
		panic(fmt.Sprintf("wrong number of params for fn %s. Expected %d, Recived %d",
			c.fnName, number, len(c.Response)))
	}

	// Validate type of params values
	for i := 0; i < number; i++ {
		expected := expectedFn(i)
		received := reflect.TypeOf(params[i])

		// Skip validation of nil values
		// TODO: For input params is ok to skip nil? Have to check if nil is acceptable for the type
		if received == nil {
			continue
		}

		// When an interfaced is expected check implementation
		if expected.Kind() == reflect.Interface && received.Implements(expected) {
			continue
		}

		// Check if the expected and received values are of the same type
		if expected.Kind() != received.Kind() {
			panic(fmt.Sprintf("func %s expect value %d to be %s, given %s",
				c.fnName, i, expected.Name(), received.Name()))
		}
	}
}

func (c *Call) Fill(params ...interface{}) {
	if len(c.Response) != len(params) {
		panic("The numbers of values doesn't match the number of params")
	}

	for i, v := range c.Response {
		// TODO: v has to be a pointer.. validate that

		// If mocked value is nil, we use the zero value of the param
		if v == nil {
			continue
		}

		param := reflect.ValueOf(params[i]).Elem()
		param.Set(reflect.ValueOf(v))
	}
}
