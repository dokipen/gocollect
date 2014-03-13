/*
github.com/dokipen/debug is a port of visionmedia/debug nodejs lib. It is meant to be extremely quick to get going with. All configuration is done on the commandline with the DEBUG environmental variable.

The current differences with nodejs's debug are that timing is not supported and we include the debug callers file and line number in the output.
*/
package debug

import (
	"crypto/md5"
	"fmt"
	colorfmt "github.com/wsxiaoys/terminal/color"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
    "time"
)

// Possible color values.
const COLORS = "rgbcmykw"

type message struct {
    from_file string
    from_line int
	namespace string
	message   string
	args      []interface{}
}

type messageChan chan *message

var (
	matchCache  = map[string]bool{}
	fmtCache    = map[string]string{}
	colorCache  = map[string]string{}
	globs       []string
    msgCh       = make(messageChan)
)


func (msg *message) log(elapsed time.Duration) {
	if match(msg.namespace) {
        payload := fmt.Sprintf(msg.message, msg.args...)
        content := fmt.Sprintf(fmtCache[msg.namespace], msg.from_file, msg.from_line, payload, elapsed)
        colorfmt.Print(content)
	}
}

func worker(ch messageChan) {
    last := time.Now()
    for {
        msg := <-ch
        elapsed := time.Since(last)
        last = time.Now()
        msg.log(elapsed)
    }
}

/*
init prepares the environmental DEBUG variable.
*/
func init() {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		globs = []string{}
	} else {
		globs = strings.Split(debug, " ")
	}
    go worker(msgCh)
}

/*
match checks if the namespace is enabled for debug logging. It also seeds the
cache for the namespace to speed things up.
*/
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

/*
getcolor returns the color of the namespace.
*/
func getcolor(namespace string) string {
	h := md5.New()
	io.WriteString(h, namespace)
	var sum int
	sumbytes := h.Sum(nil)
	for _, i := range sumbytes {
		sum += int(i)
	}
	return string(COLORS[sum%len(COLORS)])

}

/*
Log sends the log message to the log worker to be written to stdout.
*/
func Log(namespace, msg string, args ...interface{}) {
    _, path, line, _ := runtime.Caller(2)
    message := &message{
        from_file: filepath.Base(path),
        from_line: line,
        namespace: namespace,
        message: msg,
        args: args,
    }
    msgCh <- message
}

/*
Logger returns a closure for namespacing Log.
*/
func Logger(namespace string) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		Log(namespace, msg, args...)
	}
}
