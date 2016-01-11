package gocollect

import (
    "github.com/dokipen/debug"
    "strconv"
    "net"
    "http"
    "json"
    accesslog "github.com/mash/go-accesslog"
)

var (
    log = debug.Logger("gocollect")
)

func Bind(bind string) net.Listener {
    for {
        l, err := net.Listen("tcp", fmt.Sprintf(":%d", BIND))
        if err != nil {
            log("Failed to bind, retrying: %+v", err)
        } else {
            log("Listening on %s", BIND)
            return l
        }
    }
}

type logger struct {
}

func (l logger) Log(l accesslog.LogRecord) {
    log("%s - %s [%s] \"%s %s %s\" %d %d %s", l.Ip, l.Username, l.Time, l.Protocol, l.Method, l.Uri, l.Status, l.Size)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

    w.Write("{\"error\":false}")
}

func Start(listener) {
    l := logger{}
    http.HandleFunc("/", accesslog.NewLoggingHandler(Handler, l))
    http.Server(listener)
}
