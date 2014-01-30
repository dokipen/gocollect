package debug

import (
	"crypto/md5"
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const COLORS = "rgbcmykw"

var (
	matchCache  = map[string]bool{}
	prefixCache = map[string]string{}
	mutex       = new(sync.Mutex)
	globs       []string
)

func init() {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		globs = []string{}
	} else {
		globs = strings.Split(debug, " ")
	}
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

	mutex.Lock()
	defer mutex.Unlock()
	if match, ok = matchCache[namespace]; ok {
		return
	}

	for _, glob := range globs {
		if ok, _ = filepath.Match(glob, namespace); ok {
			match = true
			matchCache[namespace] = match
			prefixCache[namespace] = fmt.Sprintf("@%s%s@| ", getcolor(namespace), namespace)
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

func printns(namespace string) {
	var ok bool
	var prefix string

	mutex.Lock()
	prefix, ok = prefixCache[namespace]
	mutex.Unlock()
	if ok {
		color.Print(prefix)
	}
	return
}

func Log(namespace, msg string, args ...interface{}) {
	if match(namespace) {
		printns(namespace)
		fmt.Printf(msg, args...)
        fmt.Println("")
	}
}

func Logger(namespace string) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		Log(namespace, msg, args...)
	}
}
