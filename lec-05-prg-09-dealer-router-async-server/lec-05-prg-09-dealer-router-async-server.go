package main

import (
    "fmt"
    "log"
    "os"
    "runtime"
    "time"

    zmq "github.com/pebbe/zmq4"
)

func workerRoutine(id int, ctx *zmq.Context) {
    worker, _ := ctx.NewSocket(zmq.DEALER)
    worker.Connect("inproc://backend")
    fmt.Printf("Worker #%d started\n", id)

    for {
        frames, _ := worker.RecvMessage(0)
        ident := frames[0]
        msg := frames[1]

        fmt.Printf("Worker #%d received '%s' from %s\n", id, msg, ident)
        worker.SendMessage(ident, msg)
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    numWorkers := 0
    fmt.Sscanf(os.Args[1], "%d", &numWorkers)
    ctx, _ := zmq.NewContext()

    frontend, _ := ctx.NewSocket(zmq.ROUTER)
    frontend.Bind("tcp://*:5570")

    backend, _ := ctx.NewSocket(zmq.DEALER)
    backend.Bind("inproc://backend")

    for i := 0; i < numWorkers; i++ {
        go workerRoutine(i, ctx)
    }

    err := zmq.Proxy(frontend, backend, nil)
    frontend.Close()
    backend.Close()
    ctx.Term()
}
