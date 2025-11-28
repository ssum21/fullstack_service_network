package main

import (
    "fmt"
    "os"
    "time"

    zmq "github.com/pebbe/zmq4"
)

func main() {
	id := os.Args[1]
    identity := id
    ctx, _ := zmq.NewContext()
    socket, _ := ctx.NewSocket(zmq.DEALER)

    socket.SetIdentity(identity)
    socket.Connect("tcp://localhost:5570")

    fmt.Printf("Client %s started\n", identity)

    poller := zmq.NewPoller()
    poller.Add(socket, zmq.POLLIN)

    reqs := 0

    for {
        reqs++
        msg := fmt.Sprintf("request #%d", reqs)
        fmt.Println("Req #", reqs, "sent..")

        socket.Send(msg, 0)

        time.Sleep(1 * time.Second)

        polled, _ := poller.Poll(1 * time.Second)
        if len(polled) > 0 {
            reply, _ := socket.Recv(0)
            fmt.Printf("%s received: %s\n", identity, reply)
        }
    }
    socket.Close()
    ctx.Term()
}
