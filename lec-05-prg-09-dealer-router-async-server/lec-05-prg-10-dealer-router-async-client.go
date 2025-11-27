package main

import (
    "fmt"
    "os"
    "time"

    zmq "github.com/pebbe/zmq4"
)

func main() {
    clientID := os.Args[1]
    fmt.Println("Client", clientID, "started")

    ctx, _ := zmq.NewContext()
    socket, _ := ctx.NewSocket(zmq.DEALER)
    socket.SetIdentity(clientID)
    socket.Connect("tcp://localhost:5570")

    req := 0
    for {
        req++
        msg := fmt.Sprintf("request #%d", req)
        fmt.Println(clientID, "sent:", msg)

        socket.Send(msg, 0)
        time.Sleep(1 * time.Second)
    }
}
