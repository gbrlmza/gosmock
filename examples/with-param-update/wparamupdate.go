package wparam

import "github.com/gbrlmza/gosmock"

// The mocked implementation
type MyStruct struct {
	gosmock.MockTool
}

func (s *MyStruct) AddItem(name string, quantity int, updatableParam *string) (p1 error) {
	s.GetMockedResponse(s.AddItem, name, quantity, updatableParam).
		Update(2, &updatableParam).
		Fill(&p1)
	return
}
