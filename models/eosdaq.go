package models

import (
	"burgundy/util"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Ticker struct {
	ID              uint   `json:"id"`
	TickerName      string `json:"tickerName"`
	TokenSymbol     string `json:"tokenSymbol"`
	BaseSymbol      string `json:"baseSymbol"`
	TokenAccount    string `json:"tokenAccount"`
	ContractAccount string `json:"contractAccount"`
	CurrentPrice    int    `json:"currentPrice"`
	PrevPrice       int    `json:"prevPrice"`
	Volume          uint   `json:"volume"`
}

func TickerInit(baseSymbol string) []*Ticker {
	tickers := []*Ticker{
		&Ticker{TickerName: "Everipedia", TokenSymbol: "IQ", TokenAccount: "everipediaiq"},
		&Ticker{TickerName: "Oracle Chain", TokenSymbol: "OCT", TokenAccount: "octtothemoon"},
		&Ticker{TickerName: "Chaince", TokenSymbol: "CET", TokenAccount: "eosiochaince"},
		&Ticker{TickerName: "MEET.ONE", TokenSymbol: "MEETONE", TokenAccount: "eosiomeetone"},
		&Ticker{TickerName: "eosDAC", TokenSymbol: "EOSDAC", TokenAccount: "eosdactokens"},
		&Ticker{TickerName: "Horus Pay", TokenSymbol: "HORUS", TokenAccount: "horustokenio"},
		&Ticker{TickerName: "KARMA", TokenSymbol: "KARMA", TokenAccount: "therealkarma"},
		&Ticker{TickerName: "eosBlack", TokenSymbol: "BLACK", TokenAccount: "eosblackteam"},
		&Ticker{TickerName: "EOX Commerce", TokenSymbol: "EOX", TokenAccount: "eoxeoxeoxeox"},
		&Ticker{TickerName: "EOS Sports Bets", TokenSymbol: "ESB", TokenAccount: "esbcointoken"},
		&Ticker{TickerName: "EVR Token", TokenSymbol: "EVR", TokenAccount: "eosvrtokenss"},
		&Ticker{TickerName: "Atidium", TokenSymbol: "ATD", TokenAccount: "eosatidiumio"},
		&Ticker{TickerName: "IPOS", TokenSymbol: "IPOS", TokenAccount: "oo1122334455"},
		&Ticker{TickerName: "AdderalCoin", TokenSymbol: "ADD", TokenAccount: "eosadddddddd"},
		&Ticker{TickerName: "iRespo", TokenSymbol: "IRESPO", TokenAccount: "irespotokens"},
		&Ticker{TickerName: "Challenge DAC", TokenSymbol: "CHL", TokenAccount: "challengedac"},
		&Ticker{TickerName: "EDNA", TokenSymbol: "EDNA", TokenAccount: "ednazztokens"},
		&Ticker{TickerName: "EETH", TokenSymbol: "EETH", TokenAccount: "ethsidechain"},
		&Ticker{TickerName: "Poorman Token", TokenSymbol: "POOR", TokenAccount: "poormantoken"},
		&Ticker{TickerName: "RIDL", TokenSymbol: "RIDL", TokenAccount: "ridlridlcoin"},
		&Ticker{TickerName: "TRYBE", TokenSymbol: "TRYBE", TokenAccount: "trybenetwork"},
		&Ticker{TickerName: "WiZZ", TokenSymbol: "WIZZ", TokenAccount: "wizznetwork1"},
	}
	for i, t := range tickers {
		triBase := util.ConvertBase(i, 6)
		t.ContractAccount = strings.Replace(fmt.Sprintf("eosdaq%06s", triBase), "0", "o", -1)
		t.BaseSymbol = baseSymbol
	}
	return tickers
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
	OBID          uint   `json:"obid" gorm:"primary_key"`
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	Quantity      string `json:"quantity"`
	Volume        int
	OrderTimeJSON string    `json:"ordertime" gorm:"-"`
	OrderTime     time.Time `json:"orderTime"`
	Type          OrderType `json:"ordertype"`
}

func (ob *OrderBook) GetArgs() []interface{} {
	return []interface{}{ob.ID, ob.Name, ob.Price, ob.Quantity, ob.Volume, ob.OrderTime, ob.Type}
}

func (ob *OrderBook) UpdateDBField() {
	symbol := ""
	floatValue := 0.0
	_, err := fmt.Sscanf(ob.Quantity, "%.4f %s", &floatValue, &symbol)
	ob.Volume = int(floatValue * 10000)

	i, err := strconv.ParseInt(ob.OrderTimeJSON, 10, 64)
	if err != nil {
		mlog.Errorw("UpdateDBField", "order", ob, "err", err)
		return
	}
	ob.OrderTime = time.Unix(i, 0)
}

// EosdaqTX ...
type EosdaqTx struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Price         int       `json:"price"`
	Maker         string    `json:"maker"`
	MakerAsset    string    `json:"maker_asset"`
	Taker         string    `json:"taker"`
	TakerAsset    string    `json:"taker_asset"`
	OrderTimeJSON string    `json:"ordertime" gorm:"-"`
	OrderTime     time.Time `json:"orderTime"`
}

func (et *EosdaqTx) GetArgs() []interface{} {
	return []interface{}{et.ID, et.Price, et.Maker, et.MakerAsset, et.Taker, et.TakerAsset, et.OrderTime}
}

func (et *EosdaqTx) UpdateDBField() {
	i, err := strconv.ParseInt(et.OrderTimeJSON, 10, 64)
	if err != nil {
		mlog.Errorw("UpdateDBField", "tx", et, "err", err)
		return
	}
	et.OrderTime = time.Unix(i, 0)
}

func (et *EosdaqTx) GetVolume(tokenSymbol string) (r uint) {
	f, err := strconv.ParseFloat(strings.Replace(et.MakerAsset, " "+tokenSymbol, "", -1), 64)
	fmt.Printf("first f[%f] e[%s]\n", f, err)
	if err != nil {
		f, err = strconv.ParseFloat(strings.Replace(et.TakerAsset, " "+tokenSymbol, "", -1), 64)
		fmt.Printf("second f[%f] e[%s]\n", f, err)
		if err != nil {
			mlog.Infow("GetVolume Invalid Token", "m", et.MakerAsset, "t", et.TakerAsset, "s", tokenSymbol)
			return 0
		}
	}
	return uint(f * 10000)
}

type TxResponse []*EosdaqTx

func (tr TxResponse) GetRange(begin, end uint) (rb, re uint) {
	if begin == 0 {
		rb = tr[0].ID
	} else {
		rb = begin
	}
	if len(tr) > 1 {
		re = tr[len(tr)-2].ID
	} else {
		re = end
	}
	return
}
