package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/shirr9/order-api/internal/config"
	"log/slog"
	"net/http"
	"os"
)

func main2() {
	rt := chi.NewRouter()
	rt.Get("/order/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte("hello " + id))
	})
	http.ListenAndServe("localhost:8080", rt)
}

// SetupLogger sets up slog "dev" logger
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}

func main() {
	path := "C:\\Users\\User\\GolandProjects\\order-api\\configs\\config.yaml"
	cfg, err := config.Load(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}
