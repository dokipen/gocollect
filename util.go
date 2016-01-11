package gocollect

import (
	"os"
)

func getenvOr(name, defaultVal string) (val string) {
	val = os.Getenv(name)
	if val == "" {
		val = defaultVal
	}
	return
}
