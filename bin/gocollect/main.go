package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    gocollect "../.."
)

var (
    PID string
    BIND string
)

func init() {
    PID = gocollect.GetenvOr("PID", "")
    BIND = gocollect.GetenvOr("BIND", ":8000")

    if PID != "" {
        pid := []byte(fmt.Sprintf("%d", os.Getpid()))

        if err := ioutil.WriteFile(PID, pid, 0644); err != nil {
            panic(fmt.Sprintf("Failed to write pidfile \"%s\" with \"%+v\"\n", PID, err))
        }
    }

    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
    listener := gocollect.Bind(BIND)
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-c
        listener.Close()
    }()
    gocollect.Start(listener)
}
