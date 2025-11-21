package main

import (
	zmq "github.com/pebbe/zmq4"

    "fmt"
    "math/rand"
    "time"
)

func main() {
    sub, _ := zmq.NewSocket(zmq.SUB) // SUB
    defer sub.Close()
    sub.Connect("tcp://localhost:5557")
    sub.SetSubscribe("")

    push, _ := zmq.NewSocket(zmq.PUSH) // PUSH
    defer push.Close()
    push.Connect("tcp://localhost:5558")

    rand.Seed(time.Now().UnixNano())

    for {
        poller := zmq.NewPoller()
        poller.Add(sub, zmq.POLLIN)
        events, _ := poller.Poll(time.Millisecond * 100)

        if len(events) > 0 {
            msg, _ := sub.Recv(0)
            fmt.Println("I: received message", msg)
        } else {
            n := rand.Intn(100)
            if n < 10 {
                s := fmt.Sprintf("%d", n)
                push.Send(s, 0)
                fmt.Println("I: sending message", n)
            }
        }
    }
}
