package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("product not found")
)

type Product struct {
	ID          string
	Name        string
	Description *string
	SalePrice   *float64
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Optional[T any] struct {
	Set   bool
	Value *T
}

type CreateProductInput struct {
	Name        string
	Description *string
	SalePrice   *float64
	Price       float64
}

type PatchProductInput struct {
	Name        Optional[string]
	Description Optional[string]
	SalePrice   Optional[float64]
	Price       Optional[float64]
}

func (in PatchProductInput) HasChanges() bool {
	return in.Name.Set || in.Description.Set || in.SalePrice.Set || in.Price.Set
}

func NewProduct(id string, in CreateProductInput, now time.Time) (Product, error) {
	p := Product{
		ID:          id,
		Name:        strings.TrimSpace(in.Name),
		Description: trimPtr(in.Description),
		SalePrice:   in.SalePrice,
		Price:       in.Price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := p.Validate(); err != nil {
		return Product{}, err
	}
	return p, nil
}

func (p *Product) ApplyPatch(in PatchProductInput, now time.Time) error {
	if !in.HasChanges() {
		return ErrValidation
	}
	if in.Name.Set {
		if in.Name.Value == nil {
			return ErrValidation
		}
		p.Name = strings.TrimSpace(*in.Name.Value)
	}
	if in.Description.Set {
		p.Description = trimPtr(in.Description.Value)
	}
	if in.SalePrice.Set {
		p.SalePrice = in.SalePrice.Value
	}
	if in.Price.Set {
		if in.Price.Value == nil {
			return ErrValidation
		}
		p.Price = *in.Price.Value
	}
	p.UpdatedAt = now
	return p.Validate()
}

func (p Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrValidation
	}
	if len(p.Name) > 255 {
		return ErrValidation
	}
	if p.Price < 0 {
		return ErrValidation
	}
	if p.SalePrice != nil && *p.SalePrice < 0 {
		return ErrValidation
	}
	if p.SalePrice != nil && *p.SalePrice > p.Price {
		return ErrValidation
	}
	return nil
}

func trimPtr(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	return &s
}
