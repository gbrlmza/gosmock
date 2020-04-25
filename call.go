package gosmock

import (
	"fmt"
	"reflect"
)

type Call struct {
	mTool       *MockTool
	fn          interface{}
	fnType      reflect.Type
	fnName      string
	Count       int
	Params      []interface{}
	hasParams   bool
	paramsHash  string
	Response    []interface{}
	ParamUpdate map[int]interface{}
}

func (c *Call) WithParams(p ...interface{}) *Call {
	c.Params = p
	c.hasParams = true
	c.paramsHash = hash(p...)
	return c
}

func (c *Call) WithParamUpdate(paramPosition int, value interface{}) *Call {
	c.ParamUpdate[paramPosition] = value
	return c
}

func (c *Call) WithResponse(r ...interface{}) *Call {
	c.Response = r
	return c
}

func (c *Call) Times(n int) *Call {
	// Validate input params
	if c.Params != nil {
		c.compareTypes(c.fnType.NumIn(), c.fnType.In, c.Params, false)
	}

	// Validate output params
	c.compareTypes(c.fnType.NumOut(), c.fnType.Out, c.Response, true)

	// Validate param update
	for paramPos, _ := range c.ParamUpdate {
		if paramPos < 0 || c.fnType.NumIn() == 0 || paramPos >= c.fnType.NumIn() {
			panic(fmt.Sprintf("wrong number of param position %d for update in fn %s", paramPos, c.fnName))
		}
		//expected := c.fnType.In(paramPos)
		//received := reflect.TypeOf(paramVal)
		//c.compareType(expected, received, false)
	}

	// Add responses
	for i := 0; i < n; i++ {
		c.mTool.responses[c.fnName] = append(c.mTool.responses[c.fnName], *c)
	}
	return c
}

func (c *Call) compareTypes(number int, getExpected func(i int) reflect.Type, values []interface{}, skipNil bool) {
	// Validate number of values
	if number != len(values) {
		panic(fmt.Sprintf("wrong number of values for fn %s. Expected %d, Recived %d",
			c.fnName, number, len(c.Response)))
	}

	// Validate type of params values
	for i := 0; i < number; i++ {
		expected := getExpected(i)
		received := reflect.TypeOf(values[i])
		c.compareType(expected, received, skipNil)
	}
}

func (c *Call) compareType(expected reflect.Type, received reflect.Type, skipNil bool) {
	// Skip validation of nil values
	if received == nil && skipNil {
		return
	}

	// Check zero value
	// nil is zero value for pointers, interfaces, maps, slices, channels and function types
	if reflect.Zero(expected) == reflect.ValueOf(received) {
		return
	}

	// When an interfaced is expected check implementation
	if expected.Kind() == reflect.Interface && received.Implements(expected) {
		return
	}

	// Check if the expected and received values are of the same type
	expectedKind := expected.Kind()
	receivedKind := received.Kind()
	if expectedKind != receivedKind {
		panic(fmt.Sprintf("func %s expect value to be %s, given %s",
			c.fnName, expected.Name(), received.Name()))
	}
}

func (c *Call) Fill(params ...interface{}) {
	if len(c.Response) != len(params) {
		panic("The numbers of values doesn't match the number of params")
	}

	for i, v := range c.Response {
		// Param must be a pointer
		pType := reflect.TypeOf(params[i])
		if pType == nil || pType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("func %s expect fill param %d to be a pointer",
				c.fnName, i))
		}

		// If mocked value is nil, we use the zero value of the param
		if v == nil {
			continue
		}

		param := reflect.ValueOf(params[i]).Elem()
		param.Set(reflect.ValueOf(v))
	}
}

func (c *Call) Update(position int, param interface{}) *Call {
	// Skip if there is no value to update
	v, ok := c.ParamUpdate[position]
	if !ok {
		return c
	}

	elem := reflect.ValueOf(param).Elem()
	kingStr := elem.Kind().String()
	if elem.Kind() == reflect.Ptr {
		elem.Elem().Set(reflect.ValueOf(v))
	} else {
		elem.Set(reflect.ValueOf(v))
	}

	_ = kingStr
	return c
}
