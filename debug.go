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
	"sync"
)

// Possible color values.
const COLORS = "rgbcmykw"

var (
	matchCache  = map[string]bool{}
	prefixCache = map[string]string{}
	colorCache  = map[string]string{}
	mutex       = new(sync.Mutex)
	globs       []string
)

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

	mutex.Lock()
	defer mutex.Unlock()
	if match, ok = matchCache[namespace]; ok {
		return
	}

	for _, glob := range globs {
		if ok, _ = filepath.Match(glob, namespace); ok {
			match = true
			matchCache[namespace] = match
			colorCache[namespace] = getcolor(namespace)
			prefixCache[namespace] = colorfmt.Sprint(colorize("@%s%s@|", colorCache[namespace], namespace))
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
printNs prints the namespace and returns the namespaces color.
*/
func printNs(namespace string) (color string) {
	var ok bool
	var prefix string

	mutex.Lock()
	prefix, ok = prefixCache[namespace]
	color = colorCache[namespace]
	mutex.Unlock()
	if ok {
		fmt.Print(prefix)
	}
	return
}

/*
colorize formats the string before sending it to color.
*/
func colorize(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

/*
printCaller prints the caller's filepath and line number.
*/
func printCaller(path string, line int, color string) {
	filename := filepath.Base(path)
	colorfmt.Print(colorize("@%s%s:%d@|", color, filename, line))
}

/*
printMs prints how much time has elapsed since the last message was logged.
TODO: please implement.
*/
func printMs(color string) {
	colorfmt.Print(colorize("@%s+1ms", color))
}

/*
printMsg prints the actual log message.
*/
func printMsg(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	colored := colorize("@w%s@|", msg)
	colorfmt.Print(colored)
}

/*
Log logs your message in a namespace.
*/
func Log(namespace, msg string, args ...interface{}) {
	if match(namespace) {
		_, path, line, _ := runtime.Caller(2)

		fmt.Print("  ")
		color := printNs(namespace)

		fmt.Print(" ")
		printCaller(path, line, color)

		fmt.Print(" ")
		printMsg(msg, args...)

		//fmt.Print(" ")
		//printMs(color)

		fmt.Println()
	}
}

/*
Logger returns a closure for namespacing Log.
*/
func Logger(namespace string) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		Log(namespace, msg, args...)
	}
}
