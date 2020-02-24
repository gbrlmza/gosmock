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

func (f *MockTool) ClearMocks() {
	f.responses = nil
	f.calls = nil
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
	c := f.Mock(fn).WithParams(params...)
	f.mutex.Lock()
	defer f.mutex.Unlock()

	c.Count = len(f.calls[c.fnName]) + 1
	responses := f.responses[c.fnName]

	// Find response with params
	for i := 0; i < len(responses); i++ {
		if responses[i].hasParams && responses[i].paramsHash == c.paramsHash {
			f.extractResponse(c, i)
			return c
		}
	}

	// Find generic response (without params)
	for i := 0; i < len(responses); i++ {
		if !responses[i].hasParams {
			f.extractResponse(c, i)
			return c
		}
	}

	panic(fmt.Sprintf("no mocked response available for fn %s", c.fnName))
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

func (f *MockTool) UnusedMocks() int {
	i := 0
	for k := range f.responses {
		i += len(f.responses[k])
	}
	return i
}
