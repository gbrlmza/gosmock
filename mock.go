package gosmock

import (
	"crypto"
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

type MockTools struct {
	responses map[string][]MockCall
	calls     map[string][]MockCall
	mutex     sync.Mutex
}

type MockCall struct {
	mockTools  *MockTools
	fn         interface{}
	fnType     reflect.Type
	fnName     string
	Count      int
	Params     []interface{}
	paramsHash string
	Response   []interface{}
}

func (f *MockTools) init() {
	if f.responses == nil {
		f.responses = make(map[string][]MockCall)
		f.calls = make(map[string][]MockCall)
	}
}

func (f *MockTools) hash(objs ...interface{}) string {
	digester := crypto.MD5.New()
	for k, ob := range objs {
		fmt.Fprint(digester, k)
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	return string(digester.Sum(nil))
}

func (f *MockTools) Mock(fn interface{}) *MockCall {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.init()

	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("fn is not a fn")
	}

	mockCall := &MockCall{
		mockTools: f,
		fn:        fn,
		fnType:    fnType,
		fnName:    runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(),
	}

	return mockCall
}

func (f *MockTools) GetMockedResponse(fn interface{}, params ...interface{}) *MockCall {
	m := f.Mock(fn)

	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.init()

	m.Count = len(f.calls[m.fnName])

	responses, ok := f.responses[m.fnName]
	if !ok || m.Count >= len(responses) {
		panic(fmt.Sprintf("No mocked response available for fn %s", m.fnName))
	}

	// First try to get a mocked response for the given params
	for i := 0; i < len(responses); i++ {
		if m.paramsHash == responses[i].paramsHash {
			m.Response = responses[i].Response
			f.calls[m.fnName] = append(f.calls[m.fnName], *m)
			f.deleteResponse(m.fnName, i)
			return m
		}
	}

	// If no mocked response for the given params was found, use the next response in line
	m.Response = responses[0].Response
	f.calls[m.fnName] = append(f.calls[m.fnName], *m)
	f.deleteResponse(m.fnName, 0)
	return m
}

func (f *MockTools) deleteResponse(fnName string, index int) {
	if len(f.responses[fnName]) < index {
		panic(fmt.Sprintf("Can't delete response at index %d for fn %s", index, fnName))
	}
	f.responses[fnName] = append(f.responses[fnName][:index], f.responses[fnName][index+1:]...)
}

func (f *MockTools) GetMockedCalls(fn interface{}) []MockCall {
	fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	calls := f.calls[fnName]
	return calls
}

func (m *MockCall) WithParams(params ...interface{}) *MockCall {
	m.Params = params
	m.paramsHash = m.mockTools.hash(params)
	return m
}

func (m *MockCall) WithResponse(response ...interface{}) *MockCall {
	m.Response = response
	return m
}

func (m *MockCall) Times(number int) *MockCall {
	// Validate input params
	if m.Params != nil {
		m.compareTypes(m.fnType.NumIn(), m.fnType.In, m.Params)
	}

	// Validate output params
	m.compareTypes(m.fnType.NumOut(), m.fnType.Out, m.Response)

	for i := 0; i < number; i++ {
		m.mockTools.responses[m.fnName] = append(m.mockTools.responses[m.fnName], *m)
	}
	return m
}

func (m *MockCall) compareTypes(number int, expectedFn func(i int) reflect.Type, params []interface{}) {
	// Validate number of params
	if number != len(params) {
		panic(fmt.Sprintf("wrong number of params for fn %s. Expected %d, Recived %d",
			m.fnName, number, len(m.Response)))
	}

	// Validate type of params values
	for i := 0; i < number; i++ {
		expected := expectedFn(i)
		received := reflect.TypeOf(params[i])

		// Skip validation of nil values
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
				m.fnName, i, expected.Name(), received.Name()))
		}
	}
}

func (m *MockCall) Fill(params ...interface{}) {
	if len(m.Response) != len(params) {
		panic("The numbers of values doesn't match the number of params")
	}

	for i, v := range m.Response {
		// If mocked value is nil, we use the zero value of the param
		if v == nil {
			continue
		}

		param := reflect.ValueOf(params[i]).Elem()
		param.Set(reflect.ValueOf(v))
	}
}
