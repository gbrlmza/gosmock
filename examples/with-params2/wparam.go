package wparam

import (
	"github.com/gbrlmza/gosmock"
	"net/http"
)

// The mocked implementation
type MyStruct struct {
	gosmock.MockTool
}

type Response struct {
	*http.Response
	Err error
}

func (s *MyStruct) AddItem(name string, quantity int) (p1 *Response) {
	s.GetMockedResponse(s.AddItem, name, quantity).Fill(&p1)
	return
}
