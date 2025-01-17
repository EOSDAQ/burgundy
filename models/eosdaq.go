package models

import (
	"burgundy/util"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
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
	quantPattern, _ = regexp.Compile(`(\d+\.\d{4}) ([A-Z]+)`)
}

func parseEosFloat(str string) (uint64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return uint64(0), err
	}
	return uint64(val*100000) / 10, nil
}

// Token ...
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

// TokenInit ...
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
		//&Token{Name: "TRYBE", Symbol: "TRYBE", Account: "trybenetwork"},
		//&Token{Name: "WiZZ", Symbol: "WIZZ", Account: "wizznetwork1"},
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
	BID OrderType = iota
	ASK
	MATCH
	CANCEL
	REFUND
	IGNORE
)

// String ...
func (o OrderType) String() string {
	switch o {
	case BID:
		return "stbid"
	case ASK:
		return "stask"
	case MATCH:
		return "match"
	case CANCEL:
		return "cancel"
	case REFUND:
		return "refund"
	case IGNORE:
		return "ignore"
	default:
		return ""
	}
}

// OrderBook ...
type OrderBook struct {
	OBID        uint      `json:"obid" gorm:"primary_key"`
	ID          uint      `json:"id"`
	OrderSymbol string    `json:"orderSymbol"`
	OrderTime   time.Time `json:"orderTime"`

	*EOSData
}

// GetArgs ...
func (ob *OrderBook) GetArgs() []interface{} {
	return []interface{}{
		ob.ID,
		ob.OrderSymbol,
		ob.OrderTime,
		ob.AccountName,
		ob.Price,
		ob.Volume,
		ob.Symbol,
		ob.Type,
	}
}

// UnixTime ...
type UnixTime struct {
	time.Time
}

// UnmarshalJSON ...
func (ut *UnixTime) UnmarshalJSON(data []byte) (err error) {
	strData := strings.Trim(string(data), "\"")
	if strData == "" {
		return nil
	}

	i, err := strconv.ParseInt(fmt.Sprintf("%s000", strData), 10, 64)
	if err != nil {
		mlog.Errorw("UnmarshalJSON", "data", strData, "err", err)
		return err
	}
	ut.Time = time.Unix(0, i)
	return nil
}

// OrderData ...
type OrderData struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	Price     uint64   `json:"price"`
	Quantity  string   `json:"quantity"`
	OrderTime UnixTime `json:"ordertime"`
}

func parseQuantity(quantity string) (vol uint64, sym string, err error) {
	matched := quantPattern.FindAllStringSubmatch(quantity, -1)
	if len(matched) != 1 {
		mlog.Errorw("parseQuantity", "quantity", quantity, "err", "Invalid Quantity Format")
		return 0, "", errors.NotValidf("Invalid Quantity Format")
	}
	vol, err = parseEosFloat(matched[0][1])
	if err != nil {
		mlog.Errorw("parseQuantity", "quantity", quantity, "err", "Invalid Float Format")
		return 0, "", errors.NotValidf("Invalid Float Format")
	}
	sym = matched[0][2]

	return vol, sym, nil
}

// Parse ...
func (od *OrderData) Parse(symbol string, orderType OrderType) (r *OrderBook) {

	r = &OrderBook{
		ID:          od.ID,
		OrderSymbol: symbol,
		OrderTime:   od.OrderTime.Time,
		EOSData: &EOSData{
			AccountName: od.Name,
			Price:       od.Price,
			Type:        orderType,
		},
	}

	var err error
	r.Volume, r.Symbol, err = parseQuantity(od.Quantity)
	if err != nil {
		return nil
	}

	return r
}

// ContractData ...
type ContractData struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
}

// MarshalData ...
func (cd *ContractData) MarshalData(data interface{}) (r *EOSData) {
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
	return cd.Parse()
}

// Parse ...
func (cd *ContractData) Parse() (r *EOSData) {

	var ok bool
	if _, ok = eossysAccount[cd.From]; ok {
		return nil
	}
	if _, ok = eossysAccount[cd.To]; ok {
		return nil
	}

	memos := strings.Split(cd.Memo, "@")
	if memos[0] != MATCH.String() {
		return nil
	}

	var err error
	r = &EOSData{AccountName: cd.To, Type: MATCH}
	r.Price, err = parseEosFloat(memos[1])
	if err != nil {
		mlog.Infow("ContractData Parse", "data", cd, "err", err)
		return nil
	}

	r.Volume, r.Symbol, err = parseQuantity(cd.Quantity)
	if err != nil {
		return nil
	}

	return r
}

// EosdaqTx ...
type EosdaqTx struct {
	TXID          uint      `gorm:"primary_key"`
	ID            int64     `json:"account_action_seq"`
	OrderSymbol   string    `json:"orderSymbol"`
	OrderTime     time.Time `json:"orderTime"`
	TransactionID []byte    `json:"trx_id"`

	*EOSData
}

// EOSData ...
type EOSData struct {
	// for Backend DB
	AccountName string
	Price       uint64
	Volume      uint64
	Symbol      string
	Type        OrderType
}

// GetArgs ...
func (et *EosdaqTx) GetArgs() []interface{} {
	return []interface{}{
		et.ID,
		et.OrderSymbol,
		et.OrderTime,
		et.TransactionID,
		et.AccountName,
		et.Volume,
		et.Symbol,
		et.Type,
		et.Price,
	}
}

// GetVolume ...
func (et *EosdaqTx) GetVolume(tokenSymbol string) (r uint64) {
	if et.Symbol == tokenSymbol {
		return et.Volume
	}
	return uint64(0)
}
