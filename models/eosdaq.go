package models

import (
	"burgundy/util"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	BaseSymbol      string `json:"baseSymbol"`
	Account         string `json:"account"`
	ContractAccount string `json:"contractAccount"`
	CurrentPrice    int    `json:"currentPrice"`
	PrevPrice       int    `json:"prevPrice"`
	Volume          uint   `json:"volume"`
}

func TokenInit(baseSymbol string) []*Token {
	tokens := []*Token{
		&Token{Name: "Everipedia", Symbol: "IQ", Account: "everipediaiq"},
		&Token{Name: "Oracle Chain", Symbol: "OCT", Account: "octtothemoon"},
		&Token{Name: "Chaince", Symbol: "CET", Account: "eosiochaince"},
		&Token{Name: "MEET.ONE", Symbol: "MEETONE", Account: "eosiomeetone"},
		&Token{Name: "eosDAC", Symbol: "EOSDAC", Account: "eosdactokens"},
		&Token{Name: "Horus Pay", Symbol: "HORUS", Account: "horustokenio"},
		&Token{Name: "KARMA", Symbol: "KARMA", Account: "therealkarma"},
		&Token{Name: "eosBlack", Symbol: "BLACK", Account: "eosblackteam"},
		&Token{Name: "EOX Commerce", Symbol: "EOX", Account: "eoxeoxeoxeox"},
		&Token{Name: "EOS Sports Bets", Symbol: "ESB", Account: "esbcointoken"},
		&Token{Name: "EVR Token", Symbol: "EVR", Account: "eosvrtokenss"},
		&Token{Name: "Atidium", Symbol: "ATD", Account: "eosatidiumio"},
		&Token{Name: "IPOS", Symbol: "IPOS", Account: "oo1122334455"},
		&Token{Name: "AdderalCoin", Symbol: "ADD", Account: "eosadddddddd"},
		&Token{Name: "iRespo", Symbol: "IRESPO", Account: "irespotokens"},
		&Token{Name: "Challenge DAC", Symbol: "CHL", Account: "challengedac"},
		&Token{Name: "EDNA", Symbol: "EDNA", Account: "ednazztokens"},
		&Token{Name: "EETH", Symbol: "EETH", Account: "ethsidechain"},
		&Token{Name: "Poorman Token", Symbol: "POOR", Account: "poormantoken"},
		&Token{Name: "RIDL", Symbol: "RIDL", Account: "ridlridlcoin"},
		&Token{Name: "TRYBE", Symbol: "TRYBE", Account: "trybenetwork"},
		&Token{Name: "WiZZ", Symbol: "WIZZ", Account: "wizznetwork1"},
	}
	for i, t := range tokens {
		base := util.ConvertBase(i, 6)
		t.ContractAccount = strings.Replace(fmt.Sprintf("eosdaq%06s", base), "0", "o", -1)
		t.BaseSymbol = baseSymbol
	}
	return tokens
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
	_, err := fmt.Sscanf(ob.Quantity, "%f %s", &floatValue, &symbol)
	ob.Volume = int(floatValue * 10000)

	i, err := strconv.ParseInt(fmt.Sprintf("%s000", ob.OrderTimeJSON), 10, 64)
	if err != nil {
		mlog.Errorw("UpdateDBField", "order", ob, "err", err)
		return
	}
	ob.OrderTime = time.Unix(0, i)
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
	i, err := strconv.ParseInt(fmt.Sprintf("%s000", et.OrderTimeJSON), 10, 64)
	if err != nil {
		mlog.Errorw("UpdateDBField", "tx", et, "err", err)
		return
	}
	et.OrderTime = time.Unix(0, i)
}

func (et *EosdaqTx) GetVolume(tokenSymbol string) (r uint) {
	f, err := strconv.ParseFloat(strings.Replace(et.MakerAsset, " "+tokenSymbol, "", -1), 64)
	//fmt.Printf("first f[%f] e[%s]\n", f, err)
	if err != nil {
		f, err = strconv.ParseFloat(strings.Replace(et.TakerAsset, " "+tokenSymbol, "", -1), 64)
		//fmt.Printf("second f[%f] e[%s]\n", f, err)
		if err != nil {
			mlog.Infow("GetVolume Invalid Token", "m", et.MakerAsset, "t", et.TakerAsset, "s", tokenSymbol)
			return 0
		}
	}
	return uint(f*100000) / 10
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
