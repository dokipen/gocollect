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

// Formatted in two stages. First (color, namespace, color) which is cached,
// then (file, linenum, message, elapsetime) for the actual message.
const FMT = "  @%s%s@| @!%%s:%%d@| @w%%s@| @%s+%%s@|\n"

type message struct {
	from_file string
	from_line int
	namespace string
	message   string
	args      []interface{}
}

type messageChan chan *message

var (
	fmtCache = map[string]string{}
	globs    []string
	msgCh    = make(messageChan)
)

func (msg *message) log(elapsed time.Duration) {
	if fmtstr := formatFor(msg.namespace); fmtstr != "" {
		payload := fmt.Sprintf(msg.message, msg.args...)
		content := fmt.Sprintf(fmtstr, msg.from_file, msg.from_line, payload, elapsed)
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
func formatFor(namespace string) string {
	var ok bool

	if match, ok := fmtCache[namespace]; ok {
		return match
	}

	for _, glob := range globs {
		if ok, _ = filepath.Match(glob, namespace); ok {
			color := colorFor(namespace)
			fmtCache[namespace] = fmt.Sprintf(FMT, color, namespace, color)
			return fmtCache[namespace]
		}
	}
	fmtCache[namespace] = ""
	return fmtCache[namespace]
}

/*
getcolor returns the color of the namespace.
*/
func colorFor(namespace string) string {
	h := md5.New()
	io.WriteString(h, namespace)
	var sum int
	sumbytes := h.Sum(nil)
	for _, i := range sumbytes {
		sum += int(i)
	}
	return string(COLORS[sum%len(COLORS)])

}

// Public API
/*
Log sends the log message to the log worker to be written to stdout.
*/
func Log(namespace, msg string, args ...interface{}) {
	_, path, line, _ := runtime.Caller(2)
	message := &message{
		from_file: filepath.Base(path),
		from_line: line,
		namespace: namespace,
		message:   msg,
		args:      args,
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
