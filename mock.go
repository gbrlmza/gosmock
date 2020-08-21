package gosmock

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

var mutex sync.Mutex

type MockTool struct {
	responses map[string][]Call
	calls     map[string][]Call
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
	mutex.Lock()
	defer mutex.Unlock()
	f.init()

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("fn is not a function")
	}

	c := &Call{
		mTool:       f,
		fn:          fn,
		fnType:      fnType,
		fnName:      runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
		ParamUpdate: make(map[int]interface{}, 0),
	}

	return c
}

func (f *MockTool) GetMockedResponse(fn interface{}, params ...interface{}) *Call {
	c := f.Mock(fn).WithParams(params...)
	mutex.Lock()
	defer mutex.Unlock()

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
	c.ParamUpdate = f.responses[c.fnName][i].ParamUpdate
	f.calls[c.fnName] = append(f.calls[c.fnName], *c)
	f.responses[c.fnName] = append(f.responses[c.fnName][:i], f.responses[c.fnName][i+1:]...)
}

func (f *MockTool) GetMockedCalls(fn interface{}) []Call {
	c := f.Mock(fn)
	return f.calls[c.fnName]
}

func (f *MockTool) AllMocksUsed() bool {
	for _, r := range f.responses {
		if len(r) != 0 {
			return false
		}
	}

	return true
}
