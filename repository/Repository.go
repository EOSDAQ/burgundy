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

// InitDB ...
func InitDB(burgundy conf.ViperConfig) *gorm.DB {

	mlog, _ = util.InitLog("repository", burgundy.GetString("logmode"))

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
	dbConn.LogMode(true)
	return dbConn
}

// UserRepository ...
type UserRepository interface {
	GetByID(ctx context.Context, accountName string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Store(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, accountName string) (bool, error)
}
