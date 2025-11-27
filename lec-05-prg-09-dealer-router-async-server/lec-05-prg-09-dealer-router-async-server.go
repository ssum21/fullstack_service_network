package main

import (
    "fmt"
    "log"
    "runtime"
    "time"

    zmq "github.com/pebbe/zmq4"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    ctx, _ := zmq.NewContext()

    frontend, _ := ctx.NewSocket(zmq.ROUTER)
    frontend.Bind("tcp://*:5570")
    backend, _ := ctx.NewSocket(zmq.DEALER)
    backend.Bind("inproc://backend")

    go func() {
        for i := 0; i < 3; i++ {
            time.Sleep(300 * time.Millisecond)
			worker, _ := ctx.NewSocket(zmq.DEALER)}
    }()

    fmt.Println("proxy started")

    frontend.Close()
    backend.Close()
    ctx.Term()
}
