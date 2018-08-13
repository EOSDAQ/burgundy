package repository

import (
	models "burgundy/models"
	"burgundy/util"
	"context"
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

type testQueryFunc func(t *testing.T, m sqlmock.Sqlmock, db *gorm.DB)

func newDB() (sqlmock.Sqlmock, *gorm.DB) {
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_0")
	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}

	gormDB, gerr := gorm.Open("sqlmock", "sqlmock_db_0")
	if gerr != nil {
		log.Fatalf("can't open gorm connection: %s", err)
	}
	gormDB.LogMode(true)

	return mock, gormDB.Set("gorm:update_column", true)
}

func checkMock(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func fixedFullRe(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}

func getTestOrderBooks(n int) []models.OrderBook {
	ret := []models.OrderBook{}
	for i := 0; i < n; i++ {
		o := models.OrderBook{
			ID:        uint(i),
			Name:      fmt.Sprintf("name_%d", i),
			Price:     util.RandNum(10000),
			Quantity:  fmt.Sprintf("%d.%d ICO", util.RandNum(1000), util.RandNum(10000)),
			OrderTime: uint(time.Now().UnixNano()),
			Type:      models.ASK,
		}
		ret = append(ret, o)
	}

	return ret
}

func getRowsForOrderBooks(orderbooks []models.OrderBook) *sqlmock.Rows {
	var orderbookFieldNames = []string{"id", "name", "price", "quantity", "order_time", "order_time_readable", "ordertype"}
	rows := sqlmock.NewRows(orderbookFieldNames)
	for _, o := range orderbooks {
		rows = rows.AddRow(o.ID, o.Name, o.Price, o.Quantity, o.OrderTime, o.OrderTimeReadable, o.Type)
		fmt.Printf("rows : [%v]\n", rows)
	}
	return rows
}

func TestQueries(t *testing.T) {
	funcs := []testQueryFunc{
		testUserSelectAll,
	}
	for _, f := range funcs {
		f := f // save range var
		funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		funcName = filepath.Ext(funcName)
		funcName = strings.TrimPrefix(funcName, ".")
		t.Run(funcName, func(t *testing.T) {
			t.Parallel()
			m, db := newDB()
			defer checkMock(t, m)
			f(t, m, db)
		})
	}
}

func testUserSelectAll(t *testing.T, m sqlmock.Sqlmock, db *gorm.DB) {

	//m.ExpectExec("^CREATE TABLE ").WillReturnError(nil)
	//m.ExpectQuery(fixedFullRe("CREATE TABLE `contract_tx` (`id` int unsigned AUTO_INCREMENT,`price` int,`maker` varchar(255),`maker_asset` varchar(255),`taker` varchar(255),`taker_asset` varchar(255),`symbol` varchar(255) , PRIMARY KEY (`id`))")).WithArgs(nil).WillReturnError(nil)
	coinContract := util.RandString(12)
	repo := NewGormEosdaqRepository(db, coinContract)

	expOrderBooks := getTestOrderBooks(2)
	m.ExpectQuery(fixedFullRe(fmt.Sprintf("SELECT * FROM \"%s_order_books\"", coinContract))).
		WillReturnRows(getRowsForOrderBooks(expOrderBooks))

	orderBooks, err := repo.GetOrderBook(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, expOrderBooks, orderBooks)
}
