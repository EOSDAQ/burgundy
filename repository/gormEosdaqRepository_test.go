package repository

import (
	models "burgundy/models"
	"burgundy/util"
	"context"
	"database/sql/driver"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
)

type testQueryFunc func(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository)

func newRepo() (sqlmock.Sqlmock, EosdaqRepository) {
	contDB := fmt.Sprintf("sqlmock_db_%s", util.RandString(4))
	_, mock, err := sqlmock.NewWithDSN(contDB)
	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}

	gormDB, gerr := gorm.Open("sqlmock", contDB)
	if gerr != nil {
		log.Fatalf("can't open gorm connection: %s", err)
	}
	gormDB.LogMode(true)
	gormDB.Set("gorm:update_column", true)

	mock.ExpectExec("^CREATE TABLE ").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^CREATE TABLE ").WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewGormEosdaqRepository(gormDB, "ICO")

	return mock, repo
}

func checkMock(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func fixedFullRe(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}

func getTestTxs(n uint) []*models.EosdaqTx {
	ret := []*models.EosdaqTx{}
	for i := uint(1); i <= n; i++ {
		o := &models.EosdaqTx{
			TXID:          i,
			ID:            int64(i * uint(util.RandNum(100))),
			OrderTime:     time.Now(),
			TransactionID: []byte(util.RandString(32)),
			EOSData: &models.EOSData{
				AccountName: util.RandString(12),
				Price:       uint64(util.RandNum(1000000)),
				Volume:      uint64(util.RandNum(1000000)),
				Symbol:      "ICO",
				Type:        models.OrderType(util.RandNum(4) + 1),
			},
		}
		//fmt.Printf("getTestTxs [%v]\n", o)
		ret = append(ret, o)
	}

	return ret
}

func getRowsForTxs(txs []*models.EosdaqTx) *sqlmock.Rows {
	var txFieldNames = []string{"tx_id", "id", "order_time", "transaction_id", "account_name", "price", "volume", "symbol", "type"}
	rows := sqlmock.NewRows(txFieldNames)
	for _, t := range txs {
		rows = rows.AddRow(t.TXID, t.ID, t.OrderTime, t.TransactionID, t.AccountName, t.Price, t.Volume, t.Symbol, t.Type)
		//fmt.Printf("rows : [%v]\n", rows)
	}
	return rows
}

func getArgsForTxs(txs []*models.EosdaqTx) (ret []driver.Value) {
	for _, t := range txs {
		args := t.GetArgs()
		dvargs := make([]driver.Value, len(args))
		for i, a := range args {
			dvargs[i] = a
		}
		ret = append(ret, dvargs...)
	}
	return ret
}

func getTestOrderBooks(n uint, orderType models.OrderType) []*models.OrderBook {
	ret := []*models.OrderBook{}
	for i := uint(1); i <= n; i++ {
		o := &models.OrderBook{
			OBID:        i,
			ID:          i * uint(util.RandNum(100)),
			OrderSymbol: "ICO",
			OrderTime:   time.Now(),
			EOSData: &models.EOSData{
				AccountName: util.RandString(12),
				Price:       uint64(util.RandNum(1000000)),
				Volume:      uint64(util.RandNum(1000000)),
				Symbol:      "SYS",
				Type:        orderType,
			},
		}
		//fmt.Printf("getTestOrderBooks [%v]\n", o)
		ret = append(ret, o)
	}

	return ret
}

func getRowsForOrderBooks(orderbooks []*models.OrderBook) *sqlmock.Rows {
	var orderbookFieldNames = []string{"ob_id", "id", "order_symbol", "order_time", "account_name", "price", "volume", "symbol", "type"}
	rows := sqlmock.NewRows(orderbookFieldNames)
	for _, o := range orderbooks {
		rows = rows.AddRow(o.OBID, o.ID, o.OrderSymbol, o.OrderTime, o.AccountName, o.Price, o.Volume, o.Symbol, o.Type)
		//fmt.Printf("rows : [%v]\n", rows)
	}
	return rows
}

func getArgsForOrderBooks(obs []*models.OrderBook) (ret []driver.Value) {
	for _, o := range obs {
		args := o.GetArgs()
		dvargs := make([]driver.Value, len(args))
		for i, a := range args {
			dvargs[i] = a
		}
		ret = append(ret, dvargs...)
	}
	return ret
}

func TestQueries(t *testing.T) {
	funcs := []testQueryFunc{
		testGetTransactionByID,
		testSaveTransaction,
		testGetOrderBook,
		testSaveOrderBook,
		testDeleteOrderBook,
	}
	for _, f := range funcs {
		f := f // save range var
		funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		funcName = filepath.Ext(funcName)
		funcName = strings.TrimPrefix(funcName, ".")
		t.Run(funcName, func(t *testing.T) {
			t.Parallel()
			m, repo := newRepo()
			defer checkMock(t, m)
			f(t, m, repo)
		})
	}
}

func testGetTransactionByID(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository) {

	expTx := getTestTxs(uint(1))
	req := fmt.Sprintf(`SELECT * FROM "eosdaq_txes" WHERE (id = ? and symbol = ?) ORDER BY "eosdaq_txes"."tx_id" ASC LIMIT 1`)
	m.ExpectQuery(fixedFullRe(req)).
		WillReturnRows(getRowsForTxs(expTx))

	tx, err := repo.GetTransactionByID(context.Background(), 0)
	assert.Nil(t, err)
	assert.Equal(t, expTx[0], tx)
}

func testSaveTransaction(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository) {

	expTx := getTestTxs(uint(2))
	smt := `INSERT INTO eosdaq_txes(id, order_time, transaction_id, account_name, volume, symbol, type, price) VALUES %s`
	valueStrings := []string{}
	for range expTx {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?)")
	}
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	args := getArgsForTxs(expTx)

	m.ExpectBegin()
	m.ExpectExec(fixedFullRe(smt)).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(2, 2))
	m.ExpectCommit()

	err := repo.SaveTransaction(context.Background(), expTx)
	assert.Nil(t, err)
}

func testGetOrderBook(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository) {

	for _, orderType := range []models.OrderType{models.ASK, models.BID} {
		expOrderBooks := getTestOrderBooks(uint(2), orderType)
		m.ExpectQuery(fixedFullRe(`SELECT * FROM "order_books" WHERE (type = ? and order_symbol = ?)`)).
			WithArgs([]driver.Value{orderType, "ICO"}...).
			WillReturnRows(getRowsForOrderBooks(expOrderBooks))

		orderBooks, err := repo.GetOrderBook(context.Background(), orderType)
		assert.Nil(t, err)
		assert.Equal(t, expOrderBooks, orderBooks)
	}
}

func testSaveOrderBook(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository) {

	expOrderBooks := getTestOrderBooks(uint(2), models.ASK)
	smt := `INSERT INTO order_books(id, order_symbol, order_time, account_name, price, volume, symbol, type) VALUES %s`
	valueStrings := []string{}
	for range expOrderBooks {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?)")
	}
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	args := getArgsForOrderBooks(expOrderBooks)

	m.ExpectBegin()
	m.ExpectExec(fixedFullRe(smt)).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(2, 2))
	m.ExpectCommit()

	err := repo.SaveOrderBook(context.Background(), expOrderBooks)
	assert.Nil(t, err)
}

func testDeleteOrderBook(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository) {

	expOrderBooks := getTestOrderBooks(uint(2), models.ASK)
	smt := `DELETE FROM order_books WHERE ob_id IN (%s)`
	valueArgs := []uint{}
	for _, o := range expOrderBooks {
		valueArgs = append(valueArgs, o.OBID)
	}
	smt = fmt.Sprintf(smt, util.ArrayToString(valueArgs, ","))

	m.ExpectBegin()
	m.ExpectExec(fixedFullRe(smt)).
		//WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(2, 2))
	m.ExpectCommit()

	err := repo.DeleteOrderBook(context.Background(), expOrderBooks)
	assert.Nil(t, err)
}
