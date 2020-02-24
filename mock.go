package gosmock

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

type MockTool struct {
	responses map[string][]Call
	calls     map[string][]Call
	mutex     sync.Mutex
}

func (f *MockTool) init() {
	if f.responses == nil {
		f.responses = make(map[string][]Call)
		f.calls = make(map[string][]Call)
	}
}

func (f *MockTool) Mock(fn interface{}) *Call {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.init()

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("fn is not a function")
	}

	c := &Call{
		mTool:  f,
		fn:     fn,
		fnType: fnType,
		fnName: runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
	}

	return c
}

func (f *MockTool) GetMockedResponse(fn interface{}, params ...interface{}) *Call {
	c := f.Mock(fn).WithParams(params)
	f.mutex.Lock()
	defer f.mutex.Unlock()

	c.Count = len(f.calls[c.fnName]) + 1
	responses, ok := f.responses[c.fnName]
	if !ok || c.Count > len(responses) {
		panic(fmt.Sprintf("No mocked response available for fn %s", c.fnName))
	}

	// First try to get a mocked response for the given params
	for i := 0; i < len(responses); i++ {
		if c.paramsHash == responses[i].paramsHash {
			f.extractResponse(c, i)
			return c
		}
	}

	// If no mocked response for the given params was found, use the next response in line
	f.extractResponse(c, 0)
	return c
}

func (f *MockTool) extractResponse(c *Call, i int) {
	c.Response = f.responses[c.fnName][i].Response
	f.calls[c.fnName] = append(f.calls[c.fnName], *c)
	f.responses[c.fnName] = append(f.responses[c.fnName][:i], f.responses[c.fnName][i+1:]...)
}

func (f *MockTool) GetMockedCalls(fn interface{}) []Call {
	c := f.Mock(fn)
	return f.calls[c.fnName]
}
