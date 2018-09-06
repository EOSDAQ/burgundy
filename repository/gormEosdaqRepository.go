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
	conn   *gorm.DB
	symbol string
}

// NewGormEosdaqRepository ...
func NewGormEosdaqRepository(Conn *gorm.DB, symbol string) EosdaqRepository {
	Conn.AutoMigrate(&models.EosdaqTx{})
	Conn.AutoMigrate(&models.OrderBook{})
	return &gormEosdaqRepository{Conn, symbol}
}

func (g *gormEosdaqRepository) GetLastTransactionID(ctx context.Context) int64 {
	t := &models.EosdaqTx{}
	scope := g.conn.Select([]string{"id"}).
		Where("order_symbol = ?", g.symbol).
		Order("id desc").
		First(&t)
	if scope.Error != nil {
		mlog.Errorw("GetLastTransactionID", "symbol", g.symbol, "err", scope.Error)
		return 0
	}
	if scope.RowsAffected == 0 {
		return 0
	}
	return t.ID
}

func (g *gormEosdaqRepository) GetTransactionByID(ctx context.Context, id uint) (t *models.EosdaqTx, err error) {
	t = &models.EosdaqTx{}
	scope := g.conn.Where("id = ? and symbol = ?", id, g.symbol).First(&t)
	if scope.Error != nil {
		return nil, scope.Error
	}

	if scope.RowsAffected == 0 {
		return nil, fmt.Errorf("record not found")
	}
	return t, nil
}

func (g *gormEosdaqRepository) GetTransactions(ctx context.Context, txs []*models.EosdaqTx) (dbtxs []*models.EosdaqTx, err error) {

	if len(txs) == 0 {
		return nil, nil
	}

	valueArgs := []int64{}
	for _, t := range txs {
		valueArgs = append(valueArgs, t.ID)
	}
	scope := g.conn.Where("id in (?) and symbol = ?", valueArgs, g.symbol).Find(&dbtxs)
	if scope.Error != nil {
		return nil, scope.Error
	}

	if scope.RowsAffected == 0 {
		return nil, fmt.Errorf("record not found")
	}

	return dbtxs, nil
}

func (g *gormEosdaqRepository) SaveTransaction(ctx context.Context, txs []*models.EosdaqTx) error {

	if len(txs) == 0 {
		return nil
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, t := range txs {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, t.GetArgs()...)
	}

	smt := `INSERT INTO eosdaq_txes(id, order_symbol, order_time, transaction_id, account_name, volume, symbol, type, price) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	scope := g.conn.Begin()
	if err := scope.Exec(smt, valueArgs...).Error; err != nil {
		scope.Rollback()
		return err
	}
	scope.Commit()
	return nil
}

func (g *gormEosdaqRepository) GetOrderBook(ctx context.Context, orderType models.OrderType) (obs []*models.OrderBook, err error) {
	scope := g.conn.Where("type = ? and order_symbol = ?", orderType, g.symbol).Find(&obs)
	if scope.RowsAffected == 0 {
		//fmt.Printf("record not found")
		return nil, nil
	}
	return obs, scope.Error

}

func (g *gormEosdaqRepository) SaveOrderBook(ctx context.Context, obs []*models.OrderBook) error {

	if len(obs) == 0 {
		return nil
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, o := range obs {
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, o.GetArgs()...)
	}

	smt := `INSERT INTO order_books(id, order_symbol, order_time, account_name, price, volume, symbol, type) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	scope := g.conn.Begin()
	if err := scope.Exec(smt, valueArgs...).Error; err != nil {
		scope.Rollback()
		return err
	}
	scope.Commit()
	return nil
}

func (g *gormEosdaqRepository) DeleteOrderBook(ctx context.Context, obs []*models.OrderBook) error {

	if len(obs) == 0 {
		return nil
	}

	valueArgs := []uint{}
	for _, o := range obs {
		valueArgs = append(valueArgs, o.OBID)
	}

	smt := `DELETE FROM order_books WHERE ob_id IN (%s)`
	smt = fmt.Sprintf(smt, util.ArrayToString(valueArgs, ","))

	scope := g.conn.Begin()
	if err := scope.Exec(smt).Error; err != nil {
		scope.Rollback()
		return err
	}
	scope.Commit()
	return nil
}
