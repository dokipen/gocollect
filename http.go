package gocollect

import (
    "net"
    "net/http"
    "github.com/dokipen/debug"
)

var (
    log = debug.Logger("gocollect")
)

func Bind(bind string) net.Listener {
    for {
        l, err := net.Listen("tcp", bind)
        if err != nil {
            log("Failed to bind, retrying: %+v", err)
        } else {
            log("Listening on %s", bind)
            return l
        }
    }
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

    w.Write([]byte("{\"error\":false}"))
}

func Start(listener net.Listener) {
    http.HandleFunc("/", Handler)
    http.Serve(listener, nil)
}
