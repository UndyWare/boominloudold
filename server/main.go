package main

import (
	"fmt"
)

const (
	defaultAddress = "127.0.0.1:8500"
)

func main() {
	//addr := flag.String("addr", defaultAddress, "streamer address")
	//flag.Parse()

	//Create the other items that go into a server
	//API
	//Streamer
	//mediastore
	//liststore

	ms := MediaServer{}
	fmt.Println("server starting")
	err := ms.Start()
	if err != nil {
		fmt.Printf("server closed with err: %v", err)
	}
}
