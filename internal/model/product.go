package model

import (
	"github.com/shopspring/decimal"
)

func (p *Product) TableName() string {
	return "products"
}

type (
	Products []*Product
	Product  struct {
		ID       uint            `gorm:"primaryKey"`
		Code     string          `gorm:"uniqueIndex;not null"`
		Price    decimal.Decimal `gorm:"type:decimal(10,2);not null"`
		Variants []Variant       `gorm:"foreignKey:ProductID"`
	}

	ProductResponse struct {
		Code  string  `json:"code"`
		Price float64 `json:"price"`
	}
)

func (ps Products) ToResponse() []*ProductResponse {
	products := make([]*ProductResponse, len(ps))
	for i, p := range ps {
		products[i] = &ProductResponse{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
		}
	}

	return products
}
