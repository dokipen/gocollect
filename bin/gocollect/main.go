package main

import (
    "fmt"
    "os"
    "syscall"
    "signal"
    gocollect "../.."
)

var (
    PID string
    BIND string
)

func init() {
    PID = gocollect.getenvOr("PID", nil)
    BIND = gcollect.getenvOr("BIND", ":8000")

    if PID != nil {
        pid := []byte(fmt.Sprintf("%d", os.Getpid()))

        if err := ioutil.WriteFile(PID, pid, 0644); err != nil {
            panic(fmt.Sprintf("Failed to write pidfile \"%s\" with \"%+v\"\n", PID, err)
        }
    }
}

func main() {
    listener := gocollect.Bind(BIND)
    signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-c
        listener.Close()
    }
    gocollect.Start(listener)
}
