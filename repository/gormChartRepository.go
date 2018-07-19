package repository

import (
	models "burgundy/models"
	"context"

	"github.com/jinzhu/gorm"
)

type gormChartRepository struct {
	Conn *gorm.DB
}

// NewGormChartRepository ...
func NewGormChartRepository(Conn *gorm.DB) ChartRepository {
	Conn = Conn.AutoMigrate(&models.Chart{})
	return &gormChartRepository{Conn}
}

func (g *gormChartRepository) GetByID(ctx context.Context, id string) (chart *models.Chart, err error) {
	chart = &models.Chart{ChartID: id}

	if err := g.Conn.First(&chart, "chart_id = ?", id).Error; err != nil {
		return nil, err
	}
	return chart, nil
}
