package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/demo/product-api/internal/domain"
	"github.com/demo/product-api/internal/service"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

type fakeRepo struct {
	mu       sync.Mutex
	products map[string]domain.Product
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{products: map[string]domain.Product{}}
}

func (f *fakeRepo) Create(_ context.Context, p domain.Product) (domain.Product, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.products[p.ID] = p
	return p, nil
}

func (f *fakeRepo) GetByID(_ context.Context, id string) (domain.Product, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.products[id]
	if !ok {
		return domain.Product{}, domain.ErrNotFound
	}
	return p, nil
}

func (f *fakeRepo) Update(_ context.Context, p domain.Product) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.products[p.ID]; !ok {
		return domain.ErrNotFound
	}
	f.products[p.ID] = p
	return nil
}

func TestProductService_Create(t *testing.T) {
	repo := newFakeRepo()
	fixed := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	svc := service.NewProductService(repo).
		WithClock(func() time.Time { return fixed }).
		WithIDGenerator(func() string { return "fixed-id" })

	got, err := svc.Create(context.Background(), domain.CreateProductInput{
		Name:  "Espresso",
		Price: 55,
	})
	require.NoError(t, err)
	require.Equal(t, "fixed-id", got.ID)
	require.Equal(t, "Espresso", got.Name)
	require.Equal(t, fixed, got.CreatedAt)
}

func TestProductService_CreateValidation(t *testing.T) {
	svc := service.NewProductService(newFakeRepo())
	_, err := svc.Create(context.Background(), domain.CreateProductInput{Name: "", Price: 10})
	require.ErrorIs(t, err, domain.ErrValidation)
}

func TestProductService_PatchPartialOverwrite(t *testing.T) {
	repo := newFakeRepo()
	svc := service.NewProductService(repo).WithIDGenerator(func() string { return "p1" })
	created, err := svc.Create(context.Background(), domain.CreateProductInput{
		Name:        "Mocha",
		Description: ptr("old"),
		SalePrice:   ptr(90.0),
		Price:       120,
	})
	require.NoError(t, err)

	err = svc.Patch(context.Background(), created.ID, domain.PatchProductInput{
		Description: domain.Optional[string]{Set: true, Value: ptr("new desc")},
	})
	require.NoError(t, err)

	stored, err := repo.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, "Mocha", stored.Name)
	require.Equal(t, "new desc", *stored.Description)
	require.Equal(t, 90.0, *stored.SalePrice)
	require.Equal(t, 120.0, stored.Price)
}

func TestProductService_PatchNotFound(t *testing.T) {
	svc := service.NewProductService(newFakeRepo())
	err := svc.Patch(context.Background(), "missing", domain.PatchProductInput{
		Name: domain.Optional[string]{Set: true, Value: ptr("x")},
	})
	require.ErrorIs(t, err, domain.ErrNotFound)
}
