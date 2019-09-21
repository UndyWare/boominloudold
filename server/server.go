package main

//MediaServer media server that handles requests and streams
//musc to an endpoint
type MediaServer struct {
	//Handler   requestHandler
	requests  chan (request)
	Playlists store
	Media     store
	Streamers map[string]mediaStreamer
}

type request struct {
	Recipient string
	Msg       []byte
}

type mediaStreamer interface {
	Open() error
	Close() error
	HandleRequest(interface{})
}

//A "Store" interface for media items
//Implemented by a media storage (in a filesystem on a disk)
//	or media provider (youtube, spotify, etc.)
type store interface {
	Open(string) error
	Read(string) (interface{}, error)
	Write(string, interface{}) error
}

//Start Starts the media server listening
func (ms *MediaServer) Start() error {
	//range through the request channel
	return nil
}

//Stop Stops the media server
func (ms *MediaServer) Stop() error {
	ms.requests <- request{Recipient: "self", Msg: []byte("stop")}
	return nil
}

func (ms *MediaServer) closeStreamer(id string) error {
	return ms.Streamers[id].Close()
}

//Close Closes the media server and its connections
func (ms *MediaServer) Close() error {
	ms.Stop()
	for name := range ms.Streamers {
		ms.closeStreamer(name)
		delete(ms.Streamers, name)
	}
	return err
}

func (ms *MediaServer) connectionCount() int {
	return len(ms.Streamers)
}
