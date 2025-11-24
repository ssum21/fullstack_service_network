package main

import (
    "fmt"
    "math/rand"
    "os"
    "time"

    zmq "github.com/pebbe/zmq4"
)

func main() {
    clientID := os.Args[1]
	ctx, _ := zmq.NewContext()

    sub, _ := ctx.NewSocket(zmq.SUB) // SUB
    defer sub.Close()
    sub.Connect("tcp://localhost:5557")
    sub.SetSubscribe("")

    push, _ := ctx.NewSocket(zmq.PUSH) // PUSH
    defer push.Close()
    push.Connect("tcp://localhost:5558")

    poller := zmq.NewPoller() // Poller
    poller.Add(sub, zmq.POLLIN)

    rand.Seed(time.Now().UnixNano())

    fmt.Printf("%s started (SUB + PUSH)\n", clientID)

    for {
        socks, _ := poller.Poll(100 * time.Millisecond)

        if len(socks) > 0 {
            msg, _ := sub.Recv(0)
            fmt.Printf("%s: receive status => %s\n", clientID, msg)
        } else {
            r := rand.Intn(100)
            if r < 10 {
                time.Sleep(1 * time.Second)
                msg := fmt.Sprintf("(%s:ON)", clientID)
                push.Send(msg, 0)
                fmt.Printf("%s: send status - activated\n", clientID)
            } else if r > 90 {
                time.Sleep(1 * time.Second)
                msg := fmt.Sprintf("(%s:OFF)", clientID)
                push.Send(msg, 0)
                fmt.Printf("%s: send status - deactivated\n", clientID)
            }
        }
    }
}
