package repository

import (
	models "burgundy/models"
	"burgundy/util"
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

type gormEosdaqRepository struct {
	Conn     *gorm.DB
	Contract string
}

// NewGormEosdaqRepository ...
func NewGormEosdaqRepository(Conn *gorm.DB, contract string) EosdaqRepository {
	//Conn = Conn.Table(fmt.Sprintf("%s_tx", contract)).AutoMigrate(&models.EosdaqTx{})
	//Conn = Conn.Table(fmt.Sprintf("%s_orderbook", contract)).AutoMigrate(&models.OrderBook{})
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return fmt.Sprintf("%s_%s", contract, defaultTableName)
	}
	return &gormEosdaqRepository{Conn, contract}
}

func (g *gormEosdaqRepository) GetTransactionByID(ctx context.Context, id uint) (t *models.EosdaqTx, err error) {
	t = &models.EosdaqTx{}
	scope := g.Conn.Table(fmt.Sprintf("%s_tx", g.Contract)).Where("id = ?", id).First(&t)
	if scope.Error != nil {
		return nil, scope.Error
	}

	if scope.RowsAffected == 0 {
		return nil, fmt.Errorf("record not found")
	}
	return t, nil

}

func (g *gormEosdaqRepository) SaveTransaction(ctx context.Context, txs []*models.EosdaqTx) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, t := range txs {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, t.GetArgs()...)
	}

	smt := `INSERT INTO eosdaqtx(id, price, maker, maker_asset, taker, taker_asset, ordertime) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	tx := g.Conn.Begin()
	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (g *gormEosdaqRepository) GetOrderBook(ctx context.Context) (obs []models.OrderBook, err error) {
	session := g.Conn.New()
	scope := session.Find(&obs)
	if scope.RowsAffected == 0 {
		fmt.Printf("not found record")
		return nil, nil
	}
	/*
		for _, f := range scope.SelectAttrs() {
			fmt.Println("fields : ", f)
		}
	*/
	return obs, scope.Error

}

func (g *gormEosdaqRepository) SaveOrderBook(ctx context.Context, obs []*models.OrderBook) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, o := range obs {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, o.GetArgs()...)
	}

	smt := `INSERT INTO eosdaqtx(id, price, maker, maker_asset, taker, taker_asset, ordertime) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	tx := g.Conn.Begin()
	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (g *gormEosdaqRepository) DeleteOrderBook(ctx context.Context, obs []*models.OrderBook) error {
	valueArgs := []uint{}

	for _, o := range obs {
		valueArgs = append(valueArgs, o.ID)
	}

	smt := `DELETE FROM orderbook WHERE id IN (%s)`
	smt = fmt.Sprintf(smt, util.ArrayToString(valueArgs, ","))

	tx := g.Conn.Begin()
	if err := tx.Exec(smt).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
