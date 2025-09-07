package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shirr9/order-api/internal/order"
)

type OrderAdder interface {
	AddOrder(ctx context.Context, o order.Order) error
}

func NewAddHandler(log *slog.Logger, adder OrderAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.NewAddHandler"
		log := log.With(slog.String("op", op))

		var o order.Order
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			log.Error("failed to decode request body", slog.Any("err", err))
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := adder.AddOrder(r.Context(), o); err != nil {
			log.Error("failed to add order", slog.Any("err", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Info("order added successfully via http", slog.String("order_uid", o.OrderUID))

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "created", "order_uid": o.OrderUID})
	}
}
