package main

import (
    "fmt"
    "os"
    "runtime"

    zmq "github.com/pebbe/zmq4"
)

func ServerWorker(ctx *zmq.Context, id int) {
    worker, _ := ctx.NewSocket(zmq.DEALER)
    worker.Connect("inproc://backend")
    fmt.Printf("Worker#%d started\n", id)

    for {
        frames, _ := worker.RecvMessage(0)
        if len(frames) < 2 {
            continue
        }

        ident := frames[0]
        msg := frames[1]

        fmt.Printf("Worker#%d received %s from %s\n", id, msg, ident)
        worker.SendMessage(ident, msg)
    }

    worker.Close()
}

func ServerTask(numServer int) {
    ctx, _ := zmq.NewContext()

    frontend, _ := ctx.NewSocket(zmq.ROUTER)
    frontend.Bind("tcp://*:5570")

    backend, _ := ctx.NewSocket(zmq.DEALER)
    backend.Bind("inproc://backend")

    for i := 0; i < numServer; i++ {
        go ServerWorker(ctx, i)
    }

    zmq.Proxy(frontend, backend, nil)
    frontend.Close()
    backend.Close()
    ctx.Term()
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    var num int
    fmt.Sscanf(os.Args[1], "%d", &num)
    ServerTask(num)
}
