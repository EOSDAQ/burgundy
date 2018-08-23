package repository

import (
	"burgundy/conf"
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

type testQueryFunc func(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository, contract string)

func newRepo(contract string) (sqlmock.Sqlmock, EosdaqRepository) {
	contDB := fmt.Sprintf("sqlmock_db_%s", contract)
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
	mock.ExpectExec("^CREATE TABLE ").WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewGormEosdaqRepository(conf.Burgundy, gormDB, contract)

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

func getTestTxs(n int) []*models.EosdaqTx {
	ret := []*models.EosdaqTx{}
	for i := 1; i <= n; i++ {
		o := &models.EosdaqTx{
			ID:            uint(i),
			Price:         util.RandNum(10000),
			Maker:         util.RandString(12),
			MakerAsset:    fmt.Sprintf("%d.%d ICO", util.RandNum(1000), util.RandNum(10000)),
			Taker:         util.RandString(12),
			TakerAsset:    fmt.Sprintf("%d.%d EOS", util.RandNum(1000), util.RandNum(10000)),
			OrderTimeJSON: fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		fmt.Printf("getTestTxs [%v]\n", o)
		ret = append(ret, o)
	}

	return ret
}

func getRowsForTxs(txs []*models.EosdaqTx) *sqlmock.Rows {
	var txFieldNames = []string{"id", "price", "maker", "maker_asset", "taker", "taker_asset", "order_time"}
	rows := sqlmock.NewRows(txFieldNames)
	for _, t := range txs {
		t.UpdateDBField()
		t.OrderTimeJSON = "" // for assert
		rows = rows.AddRow(t.ID, t.Price, t.Maker, t.MakerAsset, t.Taker, t.TakerAsset, t.OrderTime)
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

func getTestOrderBooks(n int, orderType models.OrderType) []*models.OrderBook {
	ret := []*models.OrderBook{}
	for i := 1; i <= n; i++ {
		o := &models.OrderBook{
			OBID:          uint(i),
			ID:            uint(i * util.RandNum(100)),
			Name:          fmt.Sprintf("name_%d", i),
			Price:         util.RandNum(10000),
			Quantity:      fmt.Sprintf("%d.%d ICO", util.RandNum(1000), util.RandNum(10000)),
			OrderTimeJSON: fmt.Sprintf("%d", time.Now().UnixNano()),
			Type:          orderType,
		}
		fmt.Printf("getTestOrderBooks [%v]\n", o)
		ret = append(ret, o)
	}

	return ret
}

func getRowsForOrderBooks(orderbooks []*models.OrderBook) *sqlmock.Rows {
	var orderbookFieldNames = []string{"ob_id", "id", "name", "price", "quantity", "volume", "order_time", "type"}
	rows := sqlmock.NewRows(orderbookFieldNames)
	for _, o := range orderbooks {
		o.UpdateDBField()
		o.OrderTimeJSON = "" // for assert
		rows = rows.AddRow(o.OBID, o.ID, o.Name, o.Price, o.Quantity, o.Volume, o.OrderTime, o.Type)
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
			coinContract := util.RandString(12)
			m, repo := newRepo(coinContract)
			defer checkMock(t, m)
			f(t, m, repo, coinContract)
		})
	}
}

func testGetTransactionByID(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository, contract string) {

	expTx := getTestTxs(1)
	req := fmt.Sprintf(`SELECT * FROM "%s_txs" WHERE (id = ?) ORDER BY "%s_txs"."id" ASC LIMIT 1`, contract, contract)
	m.ExpectQuery(fixedFullRe(req)).
		WillReturnRows(getRowsForTxs(expTx))

	tx, err := repo.GetTransactionByID(context.Background(), 0)
	assert.Nil(t, err)
	assert.Equal(t, expTx[0], tx)
}

func testSaveTransaction(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository, contract string) {

	expTx := getTestTxs(2)
	smt := `INSERT INTO %s_txs(id, price, maker, maker_asset, taker, taker_asset, order_time) VALUES %s`
	valueStrings := []string{}
	for range expTx {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?)")
	}
	smt = fmt.Sprintf(smt, contract, strings.Join(valueStrings, ","))
	args := getArgsForTxs(expTx)

	m.ExpectBegin()
	m.ExpectExec(fixedFullRe(smt)).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(2, 2))
	m.ExpectCommit()

	err := repo.SaveTransaction(context.Background(), expTx)
	assert.Nil(t, err)
}

func testGetOrderBook(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository, contract string) {

	for _, orderType := range []models.OrderType{models.ASK, models.BID} {
		expOrderBooks := getTestOrderBooks(2, orderType)
		m.ExpectQuery(fixedFullRe(fmt.Sprintf("SELECT * FROM \"%s_order_books\" WHERE (type = ?)", contract))).
			WithArgs([]driver.Value{orderType}...).
			WillReturnRows(getRowsForOrderBooks(expOrderBooks))

		orderBooks, err := repo.GetOrderBook(context.Background(), orderType)
		assert.Nil(t, err)
		assert.Equal(t, expOrderBooks, orderBooks)
	}
}

func testSaveOrderBook(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository, contract string) {

	expOrderBooks := getTestOrderBooks(2, models.ASK)
	smt := `INSERT INTO %s_order_books(id, name, price, quantity, volume, order_time, type) VALUES %s`
	valueStrings := []string{}
	for range expOrderBooks {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?)")
	}
	smt = fmt.Sprintf(smt, contract, strings.Join(valueStrings, ","))
	args := getArgsForOrderBooks(expOrderBooks)

	m.ExpectBegin()
	m.ExpectExec(fixedFullRe(smt)).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(2, 2))
	m.ExpectCommit()

	err := repo.SaveOrderBook(context.Background(), expOrderBooks)
	assert.Nil(t, err)
}

func testDeleteOrderBook(t *testing.T, m sqlmock.Sqlmock, repo EosdaqRepository, contract string) {

	expOrderBooks := getTestOrderBooks(2, models.ASK)
	smt := `DELETE FROM %s_order_books WHERE ob_id IN (%s)`
	valueArgs := []uint{}
	for _, o := range expOrderBooks {
		valueArgs = append(valueArgs, o.OBID)
	}
	smt = fmt.Sprintf(smt, contract, util.ArrayToString(valueArgs, ","))

	m.ExpectBegin()
	m.ExpectExec(fixedFullRe(smt)).
		//WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(2, 2))
	m.ExpectCommit()

	err := repo.DeleteOrderBook(context.Background(), expOrderBooks)
	assert.Nil(t, err)
}
