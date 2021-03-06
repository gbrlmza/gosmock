package wparam

import "github.com/gbrlmza/gosmock"

// The mocked implementation
type MyStruct struct {
	gosmock.MockTool
}

func (s *MyStruct) AddItem(name string, quantity int, updatableParam interface{}) (p1 error) {
	s.GetMockedResponse(s.AddItem, name, quantity, updatableParam).
		Fill(&p1)
	return
}
