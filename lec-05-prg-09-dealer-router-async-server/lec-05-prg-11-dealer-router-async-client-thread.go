package main

import (
    "fmt"
    "os"
    "time"

    zmq "github.com/pebbe/zmq4"
)

func recvHandler(socket *zmq.Socket, identity string) {
    poller := zmq.NewPoller()
    poller.Add(socket, zmq.POLLIN)

    for {
        polled, _ := poller.Poll(1 * time.Second)
        if len(polled) > 0 {
            msg, _ := socket.Recv(0)
            fmt.Printf("%s received: %s\n", identity, msg)
        }
    }
}

func main() {
    id := os.Args[1]
    identity := id

    ctx, _ := zmq.NewContext()
    socket, _ := ctx.NewSocket(zmq.DEALER)
    socket.SetIdentity(identity)
    socket.Connect("tcp://localhost:5570")

    fmt.Printf("Client %s started\n", identity)

    go recvHandler(socket, identity)

    reqs := 0
    for {
        reqs++
        msg := fmt.Sprintf("request #%d", reqs)
        fmt.Println("Req #", reqs, "sent..")
        socket.Send(msg, 0)
        time.Sleep(1 * time.Second)
    }
    socket.Close()
    ctx.Term()
}
