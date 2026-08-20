package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/demo/product-api/internal/domain"
	httpapi "github.com/demo/product-api/internal/handler/http"
)

func TestSwaggerDocJSON(t *testing.T) {
	h := httpapi.NewProductHandler(stubService{})
	server := httptest.NewServer(httpapi.NewRouter(h))
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL + "/api-docs/doc.json")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPatchHandler_RejectsNullName(t *testing.T) {
	h := httpapi.NewProductHandler(stubService{})
	r := chi.NewRouter()
	r.Patch("/product/{id}", h.PatchProduct)

	req := httptest.NewRequest(http.MethodPatch, "/product/abc", bytes.NewBufferString(`{"name":null}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	var body httpapi.PatchAPIResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	require.Equal(t, httpapi.ErrorValidation, body.ErrorCode)
}

type stubService struct{}

func (stubService) Create(_ context.Context, _ domain.CreateProductInput) (domain.Product, error) {
	return domain.Product{}, nil
}

func (stubService) Patch(_ context.Context, _ string, _ domain.PatchProductInput) error {
	return nil
}
