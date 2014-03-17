package webizen

import (
	"log"
	"net/http"
)

type Handler struct{ http.Handler }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	log.Println(query)
	w.WriteHeader(404)
}
