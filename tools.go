package webizen

import (
	"log"
)

type URI string

func assertURI(uri URI) {
	log.Println(uri)
}
