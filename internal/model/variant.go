package model

import (
	"github.com/shopspring/decimal"
)

// Variant represents a product variant in the catalog.
// It includes a unique name, SKU, and an optional price.
// Variants can be used to represent different configurations or options for a product.

func (v *Variant) TableName() string {
	return "product_variants"
}

type (
	Variants []*Variant
	Variant  struct {
		ID        uint            `gorm:"primaryKey"`
		ProductID uint            `gorm:"not null"`
		Name      string          `gorm:"not null"`
		SKU       string          `gorm:"uniqueIndex;not null"`
		Price     decimal.Decimal `gorm:"type:decimal(10,2);null"`
	}

	VariantResponse struct {
		Name  string  `json:"name"`
		SKU   string  `json:"sku"`
		Price float64 `json:"price"`
	}
)

func (vs Variants) ToResponse(productPrice decimal.Decimal) []*VariantResponse {
	variants := make([]*VariantResponse, len(vs))
	for i, v := range vs {
		if v.Price.IsZero() {
			v.Price = productPrice
		}
		variants[i] = &VariantResponse{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.Price.InexactFloat64(),
		}
	}

	return variants
}
