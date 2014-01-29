package debug

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "github.com/wsxiaoys/terminal/color"
    "crypto/md5"
    "io"
    "sync"
)

const COLORS = "rgbcmykw"

var (
    match_cache  = map[string]bool{}
    prefix_cache = map[string]string{}
    match_mutex = new(sync.Mutex)
    prefix_mutex = new(sync.Mutex)
)

func match(namespace string) (match bool) {
    var ok bool
    if match, ok = match_cache[namespace]; !ok {
        selectors := strings.Split(os.Getenv("DEBUG"), " ")
        for selector := range selectors {
            if ok, _ = filepath.Match(selectors[selector], namespace); ok  {
                match = true
                match_mutex.Lock()
                defer match_mutex.Unlock()
                match_cache[namespace] = match
                return
            }
        }
        match = false
        match_cache[namespace] = match
    }
    return
}

func getcolor(namespace string) string {
    h := md5.New()
    io.WriteString(h, namespace)
    var sum int
    sumbytes := h.Sum(nil)
    for i := range sumbytes {
        sum += int(sumbytes[i])
    }
    return fmt.Sprintf("%c", COLORS[sum % len(COLORS)])

}

func printns(namespace string) {
    var ok bool
    var prefix string

    if prefix, ok = prefix_cache[namespace]; !ok {
        prefix = fmt.Sprintf("@%s%s@| ", getcolor(namespace), namespace)
        prefix_mutex.Lock()
        defer prefix_mutex.Unlock()
        prefix_cache[namespace] = prefix
    }

    color.Print(prefix)
}

func Log(namespace, msg string, args ...interface{}) {
    if match(namespace) {
        printns(namespace)
        fmt.Printf(msg, args...)
    }
}

func Logger(namespace string) func (string, ...interface{}) {
    return func (msg string, args ...interface{}) {
        Log(namespace, msg, args...)
    }
}
