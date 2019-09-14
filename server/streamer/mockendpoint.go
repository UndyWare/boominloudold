package streamer

import "fmt"

type mockEndpoint struct {
	Fail   bool
	IsOpen bool
	//Message []byte
	//messages chan ([]byte)
	messages chan (Request)
}

//func (me *mockEndpoint) Open(ch chan ([]byte)) error {
func (me *mockEndpoint) Open() error {
	//me.messages = ch
	if me.Fail {
		return fmt.Errorf("mock failed intentionally")
	}
	me.IsOpen = true
	return nil
}

func (me *mockEndpoint) Close() error {
	//close(me.messages)
	if me.Fail {
		return fmt.Errorf("mock failed intentionally")
	}
	me.IsOpen = false
	close(me.messages)
	return nil
}

func (me *mockEndpoint) Send(data []byte) error {
	if me.Fail {
		return fmt.Errorf("mock failed intentionally")
	}
	return nil
}

/*func (me *mockEndpoint) Requests() (func() request, error) {
	if me.Fail {
		return nil, fmt.Errorf("mock failed intentionally")
	}
	return nil, nil
}*/

func (me *mockEndpoint) Requests() chan (Request) {
	return me.messages
}
