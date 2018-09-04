package models

import (
	"burgundy/util"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	eos "github.com/eoscanada/eos-go"
)

var eossysAccount map[string]struct{}
var quantPattern *regexp.Regexp

func init() {
	eossysAccount = make(map[string]struct{})
	accounts := []string{
		"eosio.ram",
		"eosio.ramfee",
		"eosio.msig",
		"eosio.stake",
		"eosio.token",
		"eosio.saving",
		"eosio.names",
		"eosio.bpay",
		"eosio.vpay",
	}
	for _, a := range accounts {
		eossysAccount[a] = struct{}{}
	}
	quantPattern, _ = regexp.Compile("(\\d+\\.\\d{4}) ([A-Z]+)")
}

type Token struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	BaseSymbol      string `json:"baseSymbol"`
	Account         string `json:"account"`
	ContractAccount string `json:"contractAccount"`
	CurrentPrice    uint64 `json:"currentPrice"`
	PrevPrice       uint64 `json:"prevPrice"`
	Volume          uint64 `json:"volume"`
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
	NONE OrderType = iota
	ASK
	BID
	MATCH
	CANCEL
	REFUND
	IGNORE
)

// String ...
func (o OrderType) String() string {
	switch o {
	case NONE:
		return ""
	case ASK:
		return "stask"
	case BID:
		return "stbid"
	case MATCH:
		return "matched"
	case CANCEL:
		return "cancel"
	case REFUND:
		return "refund"
	default:
		return ""
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

type ContractData struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
}

func (cd *ContractData) MarshalData(token string, data interface{}) (r *EOSData) {
	src, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	keys := []string{"from", "to", "quantity", "memo"}
	for _, k := range keys {
		if _, ok := src[k]; !ok {
			return nil
		}
	}
	cd.From = src["from"].(string)
	cd.To = src["to"].(string)
	cd.Quantity = src["quantity"].(string)
	cd.Memo = src["memo"].(string)
	return cd.Parse(token)
}

func (cd *ContractData) Parse(token string) (r *EOSData) {

	r = &EOSData{}
	var err error
	matched := quantPattern.FindAllStringSubmatch(cd.Quantity, -1)
	if len(matched) != 1 {
		mlog.Infow("ContractData Parse", "data", cd, "err", "Invalid Quantity Format")
		return nil
	}
	r.Volume, err = util.ParseEosFloat(matched[0][1])
	if err != nil {
		mlog.Infow("ContractData Parse", "data", cd, "err", err)
		return nil
	}
	r.Symbol = matched[0][2]

	var ok bool
	if _, ok = eossysAccount[cd.From]; ok {
		return nil
	}
	if _, ok = eossysAccount[cd.To]; ok {
		return nil
	}

	memos := strings.Split(cd.Memo, "@")
	strPrice := ""
	switch memos[0] {
	case "match":
		r.Type = MATCH
		strPrice = memos[1]
	case "cancel":
		r.Type = CANCEL
		strPrice = memos[1]
	case "refund":
		r.Type = REFUND
		strPrice = "0.0"
	default:
		if r.Symbol == token {
			r.Type = ASK
		} else {
			r.Type = BID
		}
		strPrice = memos[0]
	}

	r.Price, err = util.ParseEosFloat(strPrice)
	if err != nil {
		mlog.Infow("ContractData Parse", "data", cd, "err", err)
		return nil
	}

	switch r.Type {
	case BID, ASK:
		r.AccountName = cd.From
	default:
		r.AccountName = cd.To
	}

	return r
}

// EosdaqTX ...
type EosdaqTx struct {
	ID            int64        `json:"account_action_seq" gorm:"primary_key"`
	OrderTime     eos.JSONTime `json:"block_time"`
	TransactionID []byte       `json:"trx_id"`

	*EOSData
}

type EOSData struct {
	// for Backend DB
	AccountName string
	Volume      uint64
	Symbol      string
	Type        OrderType
	Price       uint64
}

func (et *EosdaqTx) GetArgs() []interface{} {
	return []interface{}{
		et.ID,
		et.OrderTime,
		et.TransactionID,
		et.AccountName,
		et.Volume,
		et.Symbol,
		et.Type,
		et.Price,
	}
}

func (et *EosdaqTx) GetVolume(tokenSymbol string) (r uint64) {
	if et.Symbol == tokenSymbol {
		return et.Volume
	}
	return uint64(0)
}

type TxResponse []*EosdaqTx

func (tr TxResponse) GetRange(begin, end int64) (rb, re int64) {
	if begin == 0 {
		rb = tr[0].ID
	} else {
		rb = begin
	}
	re = tr[len(tr)-1].ID
	return
}
