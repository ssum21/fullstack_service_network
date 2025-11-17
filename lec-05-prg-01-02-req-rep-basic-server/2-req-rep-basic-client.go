package main

import (
	"context"
	"fmt"
	"github.com/go-zeromq/zmq4"
)

func main() {
	socket := zmq4.NewReq(context.Background())
	defer socket.Close()
	fmt.Println("Connecting to hello world server…")
    
	socket.Dial("tcp://localhost:5555")

	for request := 0; request < 10; request++ {
		fmt.Printf("Sending request %d …\n", request)
		socket.Send(zmq4.NewMsgString("Hello"))
		msg, _ := socket.Recv()
        fmt.Printf("Received reply %d [ %s ]\n", request, string(msg.Bytes()))
	}
}