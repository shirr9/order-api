package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// внутри создать контекст + в хэндлер логгер из мейна

func NewIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("hello " + id))
}
