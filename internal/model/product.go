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
		ID         uint            `gorm:"primaryKey"`
		Code       string          `gorm:"uniqueIndex;not null"`
		Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
		CategoryID *uint

		Category *Category `gorm:"foreignKey:CategoryID;references:ID"`
		Variants Variants  `gorm:"foreignKey:ProductID"`
	}

	ProductResponse struct {
		Code     string             `json:"code"`
		Price    float64            `json:"price"`
		Category *CategoryResponse  `json:"category"`
		Variants []*VariantResponse `json:"variants"`
	}

	ProductFilter struct {
		Category      string
		PriceLessThan decimal.Decimal
	}
)

func (ps Products) ToResponse() []*ProductResponse {
	products := make([]*ProductResponse, len(ps))
	for i, p := range ps {
		products[i] = p.ToResponse()
	}

	return products
}

func (p *Product) ToResponse() *ProductResponse {
	product := &ProductResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Variants: p.Variants.ToResponse(p.Price),
	}

	if p.Category != nil {
		product.Category = p.Category.ToResponse()
	}

	return product
}
