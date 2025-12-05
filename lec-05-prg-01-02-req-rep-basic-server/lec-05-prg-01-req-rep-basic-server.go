package main

import (
	"context"
	"fmt"
	"time"
	"github.com/go-zeromq/zmq4"
)

func main() {
	socket := zmq4.NewRep(context.Background())
	defer socket.Close()
	_ = socket.Listen("tcp://*:5555")

	for {
		msg, _ := socket.Recv()
		fmt.Printf("Received request: %s\n", string(msg.Bytes()))
		time.Sleep(1 * time.Second)
		_ = socket.Send(zmq4.NewMsgString("World"))
	}
}