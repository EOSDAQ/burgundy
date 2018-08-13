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
	Conn = Conn.Table(fmt.Sprintf("%s_eosdaq_txs", contract)).AutoMigrate(&models.EosdaqTx{})
	Conn = Conn.Table(fmt.Sprintf("%s_order_books", contract)).AutoMigrate(&models.OrderBook{})
	/*
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return fmt.Sprintf("%s_%s", contract, defaultTableName)
		}
	*/
	return &gormEosdaqRepository{Conn, contract}
}

func (g *gormEosdaqRepository) Table(table string) *gorm.DB {
	return g.Conn.Table(fmt.Sprintf("%s_%s", g.Contract, table))
}

func (g *gormEosdaqRepository) GetTransactionByID(ctx context.Context, id uint) (t *models.EosdaqTx, err error) {
	t = &models.EosdaqTx{}
	scope := g.Table("eosdaq_txs").Where("id = ?", id).First(&t)
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

	smt := `INSERT INTO %s_eosdaq_tx(id, price, maker, maker_asset, taker, taker_asset, order_time) VALUES %s`
	smt = fmt.Sprintf(smt, g.Contract, strings.Join(valueStrings, ","))

	tx := g.Conn.Begin()
	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (g *gormEosdaqRepository) GetOrderBook(ctx context.Context) (obs []*models.OrderBook, err error) {
	scope := g.Table("order_books").Find(&obs)
	if scope.RowsAffected == 0 {
		fmt.Printf("not found record")
		return nil, nil
	}
	return obs, scope.Error

}

func (g *gormEosdaqRepository) SaveOrderBook(ctx context.Context, obs []*models.OrderBook) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, o := range obs {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?)")
		valueArgs = append(valueArgs, o.GetArgs()...)
	}

	smt := `INSERT INTO %s_order_books(id, name, price, quantity, order_time, type) VALUES %s`
	smt = fmt.Sprintf(smt, g.Contract, strings.Join(valueStrings, ","))

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

	smt := `DELETE FROM %s_order_books WHERE id IN (%s)`
	smt = fmt.Sprintf(smt, g.Contract, util.ArrayToString(valueArgs, ","))

	tx := g.Conn.Begin()
	if err := tx.Exec(smt).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
