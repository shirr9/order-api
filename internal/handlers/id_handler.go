package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	errs "github.com/shirr9/order-api/internal/errors"
	"github.com/shirr9/order-api/internal/order"
)

type OrderFinder interface {
	FindOrderById(ctx context.Context, id string) (*order.Order, error)
}

func NewIdHandler(log *slog.Logger, finder OrderFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.NewIdHandler"
		log := log.With(slog.String("op", op))

		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "order_id is required", http.StatusBadRequest)
			return
		}

		orderData, err := finder.FindOrderById(r.Context(), id)
		if err != nil {
			if errors.Is(err, errs.ErrOrderNotFound) {
				http.Error(w, "order not found", http.StatusNotFound)
				return
			}
			log.Error("failed to find order", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orderData); err != nil {
			log.Error("failed to encode response", slog.Any("error", err))
		}
	}
}
