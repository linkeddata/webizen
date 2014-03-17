package webizen

import (
	"flag"
	"os"
)

var (
	bind = flag.String("bind", "", "bind address")
	dsn  = flag.String("dsn", "root@tcp(localhost:3306)/test", "")

	debug = flag.Bool("debug", false, "")
)

func init() {
	if v := os.Getenv("DSN"); len(v) > 0 {
		*dsn = v
	}
	flag.Parse()
}
