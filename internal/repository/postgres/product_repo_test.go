package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/demo/product-api/internal/domain"
	"github.com/demo/product-api/internal/repository/postgres"
	"github.com/demo/product-api/internal/testdb"
)

func ptr[T any](v T) *T { return &v }

func TestProductRepository_CreateGetUpdate(t *testing.T) {
	db := testdb.Start(t)
	repo := postgres.NewProductRepository(db)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)

	created, err := repo.Create(ctx, domain.Product{
		ID:          uuid.NewString(),
		Name:        "Matcha",
		Description: ptr("green tea"),
		SalePrice:   ptr(70.5),
		Price:       90,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "Matcha", got.Name)
	require.Equal(t, "green tea", *got.Description)
	require.InDelta(t, 70.5, *got.SalePrice, 0.001)
	require.InDelta(t, 90.0, got.Price, 0.001)

	got.Description = nil
	got.SalePrice = nil
	got.Name = "Matcha Latte"
	got.UpdatedAt = now.Add(time.Second)
	require.NoError(t, repo.Update(ctx, got))

	updated, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "Matcha Latte", updated.Name)
	require.Nil(t, updated.Description)
	require.Nil(t, updated.SalePrice)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db := testdb.Start(t)
	repo := postgres.NewProductRepository(db)
	_, err := repo.GetByID(context.Background(), uuid.NewString())
	require.ErrorIs(t, err, domain.ErrNotFound)
}
