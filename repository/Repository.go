// Package repository ...
//
// Repository will store any Database handler.
// Querying, or Creating/ Inserting into any database will stored here.
// This layer will act for CRUD to database only.
// No business process happen here. Only plain function to Database.
//
// This layer also have responsibility to choose what DB will used in Application.
// Could be Mysql, MongoDB, MariaDB, Postgresql whatever, will decided here.
//
// If using ORM, this layer will control the input, and give it directly to ORM services.
//
// If calling microservices, will handled here. Create HTTP Request to other services, and sanitize the data.
// This layer, must fully act as a repository. Handle all data input - output no specific logic happen.
//
// This Repository layer will depends to Connected DB , or other microservices if exists.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"burgundy/conf"
	models "burgundy/models"
	"burgundy/util"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql version
	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("repository", "console")
}

func makeDatabase(burgundy *conf.ViperConfig) {

	if burgundy.GetString("db_master") == "" ||
		burgundy.GetString("db_password") == "" {
		return
	}

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
		burgundy.GetString("db_master"),
		burgundy.GetString("db_password"),
		burgundy.GetString("db_host"),
		burgundy.GetInt("db_port"),
	)
	masterDB, err := sql.Open("mysql", dbURI)
	if err != nil {
		fmt.Printf("Make db open error[%s]\n", err)
		panic(err)
	}
	defer masterDB.Close()

	dbName := burgundy.GetString("db_name")
	var result string
	err = masterDB.QueryRow("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			username := burgundy.GetString("db_user")
			password := burgundy.GetString("db_pass")
			_, err = masterDB.Exec(fmt.Sprintf("CREATE DATABASE '%s'", dbName))
			if err != nil {
				panic(err)
			}
			_, err = masterDB.Exec("GRANT ALL PRIVILEGES ON '" + dbName + "'.* To '" + username + "'@'%' IDENTIFIED BY '" + password + "'")
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Make db get schema error[%s]\n", err)
			panic(err)
		}
	}
}

// InitDB ...
func InitDB(burgundy *conf.ViperConfig) *gorm.DB {

	mlog, _ = util.InitLog("repository", burgundy.GetString("logmode"))

	makeDatabase(burgundy)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		burgundy.GetString("db_user"),
		burgundy.GetString("db_pass"),
		burgundy.GetString("db_host"),
		burgundy.GetInt("db_port"),
		burgundy.GetString("db_name"),
	)
	dbConn, err := gorm.Open("mysql", dbURI) //mysql version
	if err != nil {
		fmt.Println("InitDB", err)
		os.Exit(1)
	}
	dbConn.DB().SetMaxIdleConns(100)
	if burgundy.GetString("loglevel") == "debug" {
		dbConn.LogMode(true)
	}
	return dbConn
}

// UserRepository ...
type UserRepository interface {
	GetByID(ctx context.Context, accountName string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Store(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, accountName string) (bool, error)
}

// EosdaqRepository ...
type EosdaqRepository interface {
	GetLastTransactionID(ctx context.Context) int64
	GetTransactionByID(ctx context.Context, id uint) (*models.EosdaqTx, error)
	GetTransactions(ctx context.Context, txs []*models.EosdaqTx) (dbtxs []*models.EosdaqTx, err error)
	SaveTransaction(ctx context.Context, txs []*models.EosdaqTx) error
	GetOrderBook(ctx context.Context, orderType models.OrderType) (obs []*models.OrderBook, err error)
	SaveOrderBook(ctx context.Context, obs []*models.OrderBook) error
	DeleteOrderBook(ctx context.Context, obs []*models.OrderBook) error
}

// TokenRepository ...
type TokenRepository interface {
	GetTokens(ctx context.Context) (ts []*models.Token, err error)
	GetToken(ctx context.Context, symbol string) (token *models.Token, err error)
	UpdateToken(ctx context.Context, token *models.Token) (err error)
}
