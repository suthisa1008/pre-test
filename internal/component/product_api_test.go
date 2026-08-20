package component_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	httpapi "github.com/demo/product-api/internal/handler/http"
	"github.com/demo/product-api/internal/repository/postgres"
	"github.com/demo/product-api/internal/testdb"
	"github.com/demo/product-api/internal/service"
)

func TestProductAPI_CreateAndPartialPatch(t *testing.T) {
	db := testdb.Start(t)
	repo := postgres.NewProductRepository(db)
	svc := service.NewProductService(repo)
	server := httptest.NewServer(httpapi.NewRouter(httpapi.NewProductHandler(svc)))
	t.Cleanup(server.Close)

	createBody := map[string]any{
		"name":        "Americano",
		"description": "black coffee",
		"sale_price":  45.0,
		"price":       60.0,
	}
	createResp := postJSON(t, server.URL+"/product", createBody)
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var created httpapi.APIResponse
	decode(t, createResp, &created)
	require.True(t, created.Successful)
	require.Equal(t, "", created.ErrorCode)
	require.NotNil(t, created.Data)
	require.NotEmpty(t, created.Data.Data1)
	require.Equal(t, "Americano", created.Data.Data2)

	id := created.Data.Data1
	patchResp := patchJSON(t, server.URL+"/product/"+id, map[string]any{
		"description": nil,
		"price":       65.0,
	})
	require.Equal(t, http.StatusOK, patchResp.StatusCode)

	var patched httpapi.PatchAPIResponse
	decode(t, patchResp, &patched)
	require.True(t, patched.Successful)

	stored, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "Americano", stored.Name)
	require.Nil(t, stored.Description)
	require.NotNil(t, stored.SalePrice)
	require.InDelta(t, 45.0, *stored.SalePrice, 0.001)
	require.InDelta(t, 65.0, stored.Price, 0.001)
}

func TestProductAPI_PatchNotFound(t *testing.T) {
	db := testdb.Start(t)
	server := httptest.NewServer(httpapi.NewRouter(httpapi.NewProductHandler(
		service.NewProductService(postgres.NewProductRepository(db)),
	)))
	t.Cleanup(server.Close)

	resp := patchJSON(t, server.URL+"/product/00000000-0000-0000-0000-000000000001", map[string]any{
		"name": "Nope",
	})
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var body httpapi.PatchAPIResponse
	decode(t, resp, &body)
	require.False(t, body.Successful)
	require.Equal(t, httpapi.ErrorNotFound, body.ErrorCode)
}

func postJSON(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	require.NoError(t, err)
	return resp
}

func patchJSON(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decode(t *testing.T, resp *http.Response, dest any) {
	t.Helper()
	defer resp.Body.Close()
	require.NoError(t, json.NewDecoder(resp.Body).Decode(dest))
}
