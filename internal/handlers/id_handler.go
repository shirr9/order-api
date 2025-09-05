package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("hello " + id))
}
