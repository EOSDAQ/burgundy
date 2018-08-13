package models

import (
	"fmt"
	"strconv"
	"time"
)

// Timestamp ...
type Timestamp struct {
	time.Time
}

// MarshalJSON ...
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

// UnmarshalJSON ...
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	t.Time = time.Unix(int64(ts), 0)

	return nil
}

// EosdaqTX ...
type EosdaqTx struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	Price      int       `json:"price"`
	Maker      string    `json:"maker"`
	MakerAsset string    `json:"maker_asset"`
	Taker      string    `json:"taker"`
	TakerAsset string    `json:"taker_asset"`
	OrderTime  Timestamp `json:"ordertime" gorm:"embedded"`
	Symbol     string
}

func (et *EosdaqTx) GetArgs() []interface{} {
	return []interface{}{et.ID, et.Price, et.Maker, et.MakerAsset, et.Taker, et.TakerAsset, et.OrderTime}
}

func (et *EosdaqTx) TableName() string {
	return fmt.Sprintf("%s_tx", et.Symbol)
}

// OrderType ...
type OrderType int

// OrderType types
const (
	ASK OrderType = iota
	BID
)

// String ...
func (o OrderType) String() string {
	switch o {
	case ASK:
		return "stask"
	case BID:
		return "stbid"
	default:
		return "tx"
	}
}

// OrderBook ...
type OrderBook struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Quantity  string    `json:"quantity"`
	OrderTime Timestamp `json:"ordertime"`
	Type      OrderType `json:"ordertype"`
	Symbol    string
}

func (ob *OrderBook) GetArgs() []interface{} {
	return []interface{}{ob.ID, ob.Name, ob.Price, ob.Quantity, ob.OrderTime, ob.Type}
}

func (ob *OrderBook) TableName() string {
	return fmt.Sprintf("%s_orderbook", ob.Symbol)
}
