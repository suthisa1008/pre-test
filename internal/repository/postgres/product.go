package postgres

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/demo/product-api/internal/domain"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	m := toModel(p)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return domain.Product{}, fmt.Errorf("insert product: %w", err)
	}
	return toDomain(m), nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (domain.Product, error) {
	var m ProductModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Product{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("get product: %w", err)
	}
	return toDomain(m), nil
}

func (r *ProductRepository) Update(ctx context.Context, p domain.Product) error {
	res := r.db.WithContext(ctx).Model(&ProductModel{}).Where("id = ?", p.ID).Updates(map[string]any{
		"name":        p.Name,
		"description": p.Description,
		"sale_price":  p.SalePrice,
		"price":       p.Price,
		"updated_at":  p.UpdatedAt,
	})
	if res.Error != nil {
		return fmt.Errorf("update product: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
