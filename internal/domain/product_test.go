package domain_test

import (
	"testing"
	"time"

	"github.com/demo/product-api/internal/domain"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

func TestNewProduct_Valid(t *testing.T) {
	now := time.Date(2026, 8, 20, 10, 0, 0, 0, time.UTC)
	p, err := domain.NewProduct("id-1", domain.CreateProductInput{
		Name:        " Latte ",
		Description: ptr("hot drink"),
		SalePrice:   ptr(80.0),
		Price:       100,
	}, now)
	require.NoError(t, err)
	require.Equal(t, "Latte", p.Name)
	require.Equal(t, now, p.CreatedAt)
}

func TestNewProduct_RejectsEmptyName(t *testing.T) {
	_, err := domain.NewProduct("id-1", domain.CreateProductInput{Name: "  ", Price: 10}, time.Now())
	require.ErrorIs(t, err, domain.ErrValidation)
}

func TestNewProduct_RejectsSalePriceAbovePrice(t *testing.T) {
	_, err := domain.NewProduct("id-1", domain.CreateProductInput{
		Name:      "Cake",
		SalePrice: ptr(30.0),
		Price:     20,
	}, time.Now())
	require.ErrorIs(t, err, domain.ErrValidation)
}

func TestApplyPatch_UpdatesOnlyProvidedFields(t *testing.T) {
	now := time.Now().UTC()
	p := domain.Product{
		ID:          "id-1",
		Name:        "Old",
		Description: ptr("keep me"),
		SalePrice:   ptr(10.0),
		Price:       20,
	}
	later := now.Add(time.Minute)
	err := p.ApplyPatch(domain.PatchProductInput{
		Name: domain.Optional[string]{Set: true, Value: ptr("New")},
	}, later)
	require.NoError(t, err)
	require.Equal(t, "New", p.Name)
	require.Equal(t, "keep me", *p.Description)
	require.Equal(t, 10.0, *p.SalePrice)
	require.Equal(t, 20.0, p.Price)
	require.Equal(t, later, p.UpdatedAt)
}

func TestApplyPatch_NullableDescriptionAndSalePrice(t *testing.T) {
	p := domain.Product{ID: "id-1", Name: "Old", Description: ptr("x"), SalePrice: ptr(5.0), Price: 20}
	err := p.ApplyPatch(domain.PatchProductInput{
		Description: domain.Optional[string]{Set: true, Value: nil},
		SalePrice:   domain.Optional[float64]{Set: true, Value: nil},
	}, time.Now())
	require.NoError(t, err)
	require.Nil(t, p.Description)
	require.Nil(t, p.SalePrice)
}

func TestApplyPatch_RejectsEmptyBody(t *testing.T) {
	p := domain.Product{ID: "id-1", Name: "Old", Price: 20}
	err := p.ApplyPatch(domain.PatchProductInput{}, time.Now())
	require.ErrorIs(t, err, domain.ErrValidation)
}

func TestApplyPatch_RejectsNullName(t *testing.T) {
	p := domain.Product{ID: "id-1", Name: "Old", Price: 20}
	err := p.ApplyPatch(domain.PatchProductInput{
		Name: domain.Optional[string]{Set: true, Value: nil},
	}, time.Now())
	require.ErrorIs(t, err, domain.ErrValidation)
}
