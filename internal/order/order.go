package order

import (
	"github.com/uptrace/bun"
	"time"
)

type Order struct {
	bun.BaseModel     `bun:"table:orders,alias:o"`
	OrderUID          string    `json:"order_uid" bun:"order_uid,pk" validate:"required,alphanum,len=19"`
	TrackNumber       string    `json:"track_number" bun:"track_number" validate:"required,alphanum"`
	Entry             string    `json:"entry" bun:"entry"`
	Locale            string    `json:"locale" bun:"locale" validate:"required,alpha"`
	InternalSignature string    `json:"internal_signature" bun:"internal_signature"`
	CustomerId        string    `json:"customer_id" bun:"customer_id" validate:"required,alphanum"`
	DeliveryService   string    `json:"delivery_service" bun:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" bun:"shardkey" validate:"required"`
	SmId              int       `json:"sm_id" bun:"sm_id" validate:"required,gte=0"`
	DateCreated       time.Time `json:"date_created" bun:"date_created,type:timestamp" validate:"required"`
	OofShard          string    `json:"oof_shard" bun:"oof_shard" validate:"required"`

	Delivery *Delivery `json:"delivery" bun:"rel:has-one,join:order_uid=order_uid" validate:"required"`
	Payment  *Payment  `json:"payment" bun:"rel:has-one,join:order_uid=order_uid" validate:"required"`
	Items    []Item    `json:"items" bun:"rel:has-many,join:order_uid=order_uid" validate:"required, dive"`
}

type Delivery struct {
	bun.BaseModel `bun:"table:delivery,alias:d"`
	OrderUID      string `bun:"order_uid,pk" json:"order_uid" validate:"required,alphanum,len=19"`
	Name          string `bun:"name" json:"name" validate:"required"`
	Phone         string `bun:"phone" json:"phone" validate:"required,e164"`
	Zip           string `bun:"zip" json:"zip" validate:"required,alphanum"`
	City          string `bun:"city" json:"city" validate:"required,alpha"`
	Address       string `bun:"address" json:"address" validate:"required"`
	Region        string `bun:"region" json:"region" validate:"required,alpha"`
	Email         string `bun:"email" json:"email" validate:"required,email"`
}

type Payment struct {
	bun.BaseModel `bun:"table:payment,alias:p"`
	Transaction   string `bun:"transaction,pk" json:"transaction" validate:"required"`
	OrderUID      string `bun:"order_uid" json:"order_uid" validate:"required,alphanum,len=19"`
	RequestId     string `bun:"request_id" json:"request_id"`
	Currency      string `bun:"currency" json:"currency" validate:"required,alpha,len=3"`
	Provider      string `bun:"provider" json:"provider" validate:"required"`
	Amount        int    `bun:"amount" json:"amount" validate:"required,gte=0"`
	PaymentDt     int64  `bun:"payment_dt,type:bigint" json:"payment_dt" validate:"required,gte=0"`
	Bank          string `bun:"bank" json:"bank" validate:"required"`
	DeliveryCost  int    `bun:"delivery_cost" json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal    int    `bun:"goods_total" json:"goods_total" validate:"required,gte=0"`
	CustomFee     int    `bun:"custom_fee" json:"custom_fee" validate:"required,gte=0"`
}

type Item struct {
	bun.BaseModel `bun:"table:items,alias:i"`
	ID            int64  `bun:"id,pk,autoincrement" json:"id" validate:"required,gte=0"`
	OrderUID      string `bun:"order_uid" json:"order_uid" validate:"required,alphanum,len=19"`
	ChrtId        int    `bun:"chrt_id" json:"chrt_id" validate:"required,gte=0"`
	TrackNumber   string `bun:"track_number" json:"track_number" validate:"required,alphanum"`
	Price         int    `bun:"price" json:"price" validate:"required,gte=0"`
	Rid           string `bun:"rid" json:"rid" validate:"required"`
	Name          string `bun:"name" json:"name" validate:"required"`
	Sale          int    `bun:"sale" json:"sale" validate:"gte=0,lte=100"`
	Size          string `bun:"size" json:"size" validate:"required,alphanum"`
	TotalPrice    int    `bun:"total_price" json:"total_price" validate:"required,gte=0"`
	NmId          int    `bun:"nm_id" json:"nm_id" validate:"required,gte=0"`
	Brand         string `bun:"brand" json:"brand" validate:"required"`
	Status        int    `bun:"status" json:"status" validate:"required"`
}
