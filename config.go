package webizen

import (
	"flag"
	"os"
)

var (
	debug = flag.Bool("debug", false, "")
	dsn   = flag.String("dsn", "root@tcp(localhost:3306)/test", "")
)

func init() {
	if v := os.Getenv("DSN"); len(v) > 0 {
		*dsn = v
	}
}
