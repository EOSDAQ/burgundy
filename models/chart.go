package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Chart ...
type Chart struct {
	gorm.Model `json:"-"`
	ChartID    string `json:"chartID"`
}

func (c *Chart) String() string {
	return fmt.Sprintf("chartID[%s]", c.ChartID)
}
