package streamer

import (
	"reflect"
	"testing"
)

// tt := []struct {
// 	names [string]
// 	eps [endpoints]
// 	count int
// }{
// {["conn1"], [mockEndpoint{Open:true,RW:io.ReadWriter}], 1},
// {["conn2"], [mockEndpoint{Open:true,RW:io.ReadWriter}], 1},
// {["conn3"], [mockEndpoint{Open:true,RW:io.ReadWriter}], 1},
// {["conn4"], [mockEndpoint{Open:true,RW:io.ReadWriter}], 1},
//}
// 	for i, test := range tt {
// 		if err = validateParams(test.F, test.S); (err != nil) == test.IsNil {
// 			t.Fatalf("Test (%d) failed.", i)
// 		}
// 	}

func TestNewStreamer(t *testing.T) {
	id := "127.0.0.1:8000"
	size := 15
	qu := make([]string, size)
	me := &mockEndpoint{
		Fail:   true,
		IsOpen: false,
	}

	ss := &Streamer{
		id:      id,
		ep:      me,
		playing: false,
		queue:   qu,
	}
	ns, err := NewStreamer(id, me, size)
	if err != nil {
		t.Errorf("newstreamer err: %v", err)
	}
	if !reflect.DeepEqual(*ss, *ns) {
		t.Errorf("streamers are not equal")
		t.Errorf("got: %v", *ss)
		t.Errorf("want: %v", *ns)
	}
}

func TestClose(t *testing.T) {

}
