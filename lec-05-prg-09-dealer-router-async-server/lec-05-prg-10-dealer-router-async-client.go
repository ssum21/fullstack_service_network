package main

import (
    "fmt"
    "os"
    "time"

    zmq "github.com/pebbe/zmq4"
)

func main() {
	clientID := os.Args[1]
    ctx, _ := zmq.NewContext()
    socket, _ := ctx.NewSocket(zmq.DEALER)
    socket.SetIdentity(clientID)
    socket.Connect("tcp://localhost:5570")

    poller := zmq.NewPoller()
    poller.Add(socket, zmq.POLLIN)

    req := 0
    for {
        req++
        msg := fmt.Sprintf("request #%d", req)
        fmt.Println(clientID, "sent:", msg)
        socket.Send(msg, 0)

        polled, _ := poller.Poll(1 * time.Second)
        if len(polled) > 0 {
            reply, _ := socket.Recv(0)
            fmt.Printf("%s received: %s\n", clientID, reply)
        } else {
            fmt.Println(clientID, "no reply (timeout)")
        }

        time.Sleep(1 * time.Second)
    }
}
