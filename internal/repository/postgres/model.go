package postgres

import (
	"time"

	"github.com/demo/product-api/internal/domain"
)

type ProductModel struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name;type:text;not null;index"`
	Description *string   `gorm:"column:description;type:text"`
	SalePrice   *float64  `gorm:"column:sale_price;type:double precision"`
	Price       float64   `gorm:"column:price;type:double precision;not null"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (ProductModel) TableName() string {
	return "products"
}

func toModel(p domain.Product) ProductModel {
	return ProductModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		SalePrice:   p.SalePrice,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toDomain(m ProductModel) domain.Product {
	return domain.Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		SalePrice:   m.SalePrice,
		Price:       m.Price,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
