package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/fcgi"

	"github.com/linkeddata/webizen"
)

var (
	bind = flag.String("bind", "", "bind address (empty: fcgi)")

	handler = new(webizen.Handler)
)

func init() {
	flag.Parse()
}

func main() {
	var err error

	if bind == nil || len(*bind) == 0 {
		err = fcgi.Serve(nil, handler)
	} else {
		err = http.ListenAndServe(*bind, handler)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
