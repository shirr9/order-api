package order

import (
	"github.com/uptrace/bun"
	"time"
)

type Order struct {
	bun.BaseModel     `bun:"table:orders,alias:o"`
	OrderUID          string    `json:"order_uid" bun:"order_uid,pk"`
	TrackNumber       string    `json:"track_number" bun:"track_number"`
	Entry             string    `json:"entry" bun:"entry"`
	Locale            string    `json:"locale" bun:"locale"`
	InternalSignature string    `json:"internal_signature" bun:"internal_signature"`
	CustomerId        string    `json:"customer_id" bun:"customer_id"`
	DeliveryService   string    `json:"delivery_service" bun:"delivery_service"`
	ShardKey          string    `json:"shardkey" bun:"shardkey"`
	SmId              int       `json:"sm_id" bun:"sm_id"`
	DateCreated       time.Time `json:"date_created" bun:"date_created,type:timestamp"`
	OofShard          string    `json:"oof_shard" bun:"oof_shard"`

	Delivery *Delivery `json:"delivery" bun:"rel:has-one,join:order_uid=order_uid"`
	Payment  *Payment  `json:"payment" bun:"rel:has-one,join:order_uid=order_uid"`
	Items    []Item    `json:"items" bun:"rel:has-many,join:order_uid=order_uid"`
}

type Delivery struct {
	bun.BaseModel `bun:"table:delivery,alias:d"`
	OrderUID      string `bun:"order_uid,pk" json:"order_uid"`
	Name          string `bun:"name" json:"name"`
	Phone         string `bun:"phone" json:"phone"`
	Zip           string `bun:"zip" json:"zip"`
	City          string `bun:"city" json:"city"`
	Address       string `bun:"address" json:"address"`
	Region        string `bun:"region" json:"region"`
	Email         string `bun:"email" json:"email"`
}

type Payment struct {
	bun.BaseModel `bun:"table:payment,alias:p"`
	Transaction   string `bun:"transaction,pk" json:"transaction"`
	OrderUID      string `bun:"order_uid" json:"order_uid"`
	RequestId     string `bun:"request_id" json:"request_id"`
	Currency      string `bun:"currency" json:"currency"`
	Provider      string `bun:"provider" json:"provider"`
	Amount        int    `bun:"amount" json:"amount"`
	PaymentDt     int64  `bun:"payment_dt,type:bigint" json:"payment_dt"`
	Bank          string `bun:"bank" json:"bank"`
	DeliveryCost  int    `bun:"delivery_cost" json:"delivery_cost"`
	GoodsTotal    int    `bun:"goods_total" json:"goods_total"`
	CustomFee     int    `bun:"custom_fee" json:"custom_fee"`
}

type Item struct {
	bun.BaseModel `bun:"table:items,alias:i"`
	ID            int64  `bun:"id,pk,autoincrement" json:"id"`
	OrderUID      string `bun:"order_uid" json:"order_uid"`
	ChrtId        int    `bun:"chrt_id" json:"chrt_id"`
	TrackNumber   string `bun:"track_number" json:"track_number"`
	Price         int    `bun:"price" json:"price"`
	Rid           string `bun:"rid" json:"rid"`
	Name          string `bun:"name" json:"name"`
	Sale          int    `bun:"sale" json:"sale"`
	Size          string `bun:"size" json:"size"`
	TotalPrice    int    `bun:"total_price" json:"total_price"`
	NmId          int    `bun:"nm_id" json:"nm_id"`
	Brand         string `bun:"brand" json:"brand"`
	Status        int    `bun:"status" json:"status"`
}
