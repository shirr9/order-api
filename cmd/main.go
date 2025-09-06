package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/shirr9/order-api/internal/config"
	"github.com/shirr9/order-api/internal/order"
	"github.com/shirr9/order-api/internal/storage/postgresql"
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
	//path := "C:\\Users\\User\\GolandProjects\\order-api\\configs\\config.yaml"
	//cfg, err := config.Load(path)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////ctx := context.Background()
	////st, _ := postgresql.New(ctx, cfg)
	////repo := st.NewPostgresRepository()
	path := "C:\\Users\\User\\GolandProjects\\order-api\\configs\\config.yaml"

	cfg, err := config.Load(path)
	if err != nil {
		// критично: без конфига дальше нельзя
		panic(fmt.Errorf("config load failed: %w", err))
	}

	ctx := context.Background()
	st, err := postgresql.New(ctx, cfg)
	if err != nil {
		// критично: Storage не создан, st == nil
		panic(fmt.Errorf("storage init failed: %w", err))
	}
	defer st.Close()

	repo := st.NewPostgresRepository()
	data := []byte(`{
   "order_uid": "b563feb7b2b84b6test",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}`)
	var o order.Order
	if err := json.Unmarshal(data, &o); err != nil {
		panic(fmt.Errorf("unmarshal failed: %w", err))
	}
	//// тест добавления - все норм
	//if err := repo.AddOrder(ctx, o); err != nil {
	//	panic(fmt.Errorf("AddOrder failed: %w", err))
	//}

	// тест поиска
	got, err := repo.FindOrderById(ctx, o.OrderUID)
	if err != nil {
		panic(fmt.Errorf("FindOrderById failed: %w", err))
	}
	b, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		panic(fmt.Errorf("marshal result failed: %w", err))
	}
	fmt.Println(string(b))
}
