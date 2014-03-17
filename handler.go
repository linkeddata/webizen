package webizen

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	methodsAll = []string{
		"OPTIONS", "HEAD", "GET", "POST",
	}
)

type Handler struct{ http.Handler }

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	origin := ""
	origins := req.Header["Origin"] // all CORS requests
	if len(origins) > 0 {
		origin = origins[0]
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	switch req.Method {
	case "OPTIONS":
		w.Header().Set("Accept-Patch", "application/json")
		w.Header().Set("Accept-Post", "text/turtle,application/json")

		// TODO: WAC
		corsReqH := req.Header["Access-Control-Request-Headers"] // CORS preflight only
		if len(corsReqH) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsReqH, ", "))
		}
		corsReqM := req.Header["Access-Control-Request-Method"] // CORS preflight only
		if len(corsReqM) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsReqM, ", "))
		} else {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methodsAll, ", "))
		}
		if len(origin) < 1 {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Allow", strings.Join(methodsAll, ", "))
		w.WriteHeader(200)
		return

	case "GET", "POST":
		query := req.FormValue("q")
		for _, elt := range strings.Split(query, " ") {
			if len(elt) < 6 {
				continue
			}
			if elt[:6] == "https:" || elt[:5] == "http:" {
				assertURI(URI(elt))
			}
		}
		log.Println("TODO", query)

	default:
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method Not Allowed:", req.Method)
		return
	}

	w.WriteHeader(404)
}
