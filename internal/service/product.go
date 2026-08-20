package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/demo/product-api/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	GetByID(ctx context.Context, id string) (domain.Product, error)
	Update(ctx context.Context, p domain.Product) error
}

type Clock func() time.Time

type IDGenerator func() string

type ProductService struct {
	repo  ProductRepository
	now   Clock
	newID IDGenerator
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
		now:  time.Now,
		newID: func() string {
			return uuid.NewString()
		},
	}
}

func (s *ProductService) WithClock(now Clock) *ProductService {
	s.now = now
	return s
}

func (s *ProductService) WithIDGenerator(newID IDGenerator) *ProductService {
	s.newID = newID
	return s
}

func (s *ProductService) Create(ctx context.Context, in domain.CreateProductInput) (domain.Product, error) {
	p, err := domain.NewProduct(s.newID(), in, s.now().UTC())
	if err != nil {
		return domain.Product{}, err
	}
	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}
	return created, nil
}

func (s *ProductService) Patch(ctx context.Context, id string, in domain.PatchProductInput) error {
	if id == "" {
		return domain.ErrValidation
	}
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := p.ApplyPatch(in, s.now().UTC()); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return fmt.Errorf("patch product: %w", err)
	}
	return nil
}
