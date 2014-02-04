## Debug ##
Debug is a golang logger that is meant to be extremely simple to get started
with and configure. Debug uses a single environmental variable for
configuration and loggers are initiated in one line of code.

To enable output of your logger, simply append your loggers name to the DEBUG
environmental variable. DEBUG should be a space delimited list of logger
names. Globbing is supported. Use DEBUG="\*" to enable all loggers.

## Usage ##
```go
// myprog.go
package main

import (
  "github.com/dokipen/debug"
)

var (
  log = debug.Logger("main")
)

func main() {
    log("Hello debug!")
}
```

```bash
$ DEBUG="main" go run myprog.go
main Hello debug!
```
