package main

import (
	zmq "github.com/pebbe/zmq4"

    "fmt"
)

func main() {
    ctx, _ := zmq.NewContext()

    publisher, _ := ctx.NewSocket(zmq.PUB) // PUB
    defer publisher.Close()
    publisher.Bind("tcp://*:5557")

    collector, _ := ctx.NewSocket(zmq.PULL) // PULL
    defer collector.Close()
    collector.Bind("tcp://*:5558")

    for {
        msg, _ := collector.RecvBytes(0)
        fmt.Println("I: publishing update", string(msg))
        publisher.SendBytes(msg, 0)
    }
}
