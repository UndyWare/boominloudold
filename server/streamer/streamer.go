package streamer

import (
	"encoding/json"
	"fmt"
)

//Streamer struct that handles communication between endpoints and the server
//Spawns endpoint connections
type Streamer struct {
	id string
	ep endpoint

	requests chan (Request)
	queue    []string
	playing  bool
}

//The interface for every endpoint to implement (discord, music player, etc.)
type endpoint interface {
	Open() error
	Close() error
	Send([]byte) error
	Requests() chan (Request)
	//Requests() (func() Request, error)
}

//Request Streamer request structure
type Request struct {
	Cmd string `json:"cmd"`
	Msg []byte `json:"msg"`
}

type addcmd struct {
	Name string `json:"name"`
	Typ  string `json:"type"`
	Addr string `json:"addr"`
}

//NewStreamer Creates new socket streamer from address
//TODO: see about how to return this properly
//check to see if returning an interface is the answer
func NewStreamer(id string, ep endpoint, size int) (*Streamer, error) {
	//create the new socket
	if id == "" {
		return nil, fmt.Errorf("empty addr string given for socket")
	}
	if ep == nil {
		return nil, fmt.Errorf("invalid endpoint")
	}
	ss := &Streamer{
		id:      id,
		playing: false,
		ep:      ep,
		queue:   make([]string, size),
	}
	return ss, nil
}

//HandleRequest handles incoming requests and dispatches commands
func (s *Streamer) HandleRequest(req *Request) error {
	switch req.Cmd {
	case "add":
		ac := &addcmd{}
		err := json.Unmarshal(req.Msg, ac)
		if err != nil {
			return fmt.Errorf("failed unmarshaling addcmd: %v", err)
		}
	//Playback control
	case "play":
		s.ep.
	case "pause":
	case "skip":
	case "stop": //stop playback and removes queue
	case "volume":
	//List control
	case "load": //loads given playlist
	case "display":
	case "shuffle":
	case "queue":
	case "enqueue":
	case "remove":
	//Streamer info
	case "find":
	case "quit":
	case "help":
	default:
		//log error here
	}
	return nil
}

//Close Closes the streamer and its endpoint
func (s *Streamer) Close() error {
	//TODO: how to delete a list
	s.requests <- Request{"close", nil}
	return s.ep.Close()
}
