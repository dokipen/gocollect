package debug

import (
	"crypto/md5"
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
    "time"
)

const COLORS = "rgbcmykw"

var (
	matchCache  = map[string]bool{}
	fmtCache    = map[string]string{}
	globs       []string
    msgCh       = make(chan *Message)
)

type Message struct {
    from_file string
    from_line int
	namespace string
	message   string
	args      []interface{}
}

type MessageChan chan *Message

func (msg *Message) Log(elapsed time.Duration) {
	if match(msg.namespace) {
        message := fmt.Sprintf(msg.message, msg.args...)
        content := fmt.Sprintf(fmtCache[msg.namespace], msg.from_file, msg.from_line, message, elapsed)
        color.Print(content)
	}
}

func worker(ch MessageChan) {
    last := time.Now()
    for {
        msg := <-ch
        elapsed := time.Since(last)
        last = time.Now()
        msg.Log(elapsed)
    }
}

func init() {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		globs = []string{}
	} else {
		globs = strings.Split(debug, " ")
	}
    go worker(msgCh)
}

func match(namespace string) (match bool) {
	var ok bool

	match = false

	// This might be a pointless opimization, but I'm hoping it helps with
	// multithreaded applications because we can skip the mutex if there
	// is no DEBUG env.
	if len(globs) == 0 {
		return
	}

	if match, ok = matchCache[namespace]; ok {
		return
	}

	for _, glob := range globs {
		if ok, _ = filepath.Match(glob, namespace); ok {
			match = true
			matchCache[namespace] = match
            color := getcolor(namespace)
            fmtCache[namespace] = fmt.Sprintf("  @%s%s@| @!%%s:%%d@| @w%%s@| @%s+%%s@|\n", color, namespace, color)
			return
		}
	}
	matchCache[namespace] = match
	return
}

func getcolor(namespace string) string {
	h := md5.New()
	io.WriteString(h, namespace)
	var sum int
	sumbytes := h.Sum(nil)
	for _, i := range sumbytes {
		sum += int(i)
	}
	return fmt.Sprintf("%c", COLORS[sum%len(COLORS)])

}

func Log(namespace, msg string, args ...interface{}) {
    _, path, line, _ := runtime.Caller(2)
    message := &Message{
        from_file: filepath.Base(path),
        from_line: line,
        namespace: namespace,
        message: msg,
        args: args,
    }
    msgCh <- message
}

func Logger(namespace string) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		Log(namespace, msg, args...)
	}
}
