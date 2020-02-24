package simple

import "github.com/gbrlmza/gosmock"

// The mocked implementation
type Simple struct {
	gosmock.MockTool
}

func (s *Simple) AddItem(name string, quantity int) (p1 string, p2 error) {
	s.GetMockedResponse(s.AddItem, name, quantity).Fill(&p1, &p2)
	return
}
