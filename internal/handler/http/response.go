package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/demo/product-api/internal/domain"
)

const (
	ErrorSuccess    = "SUCCESS"
	ErrorValidation = "VALIDATION_ERROR"
	ErrorNotFound   = "NOT_FOUND"
	ErrorInternal   = "INTERNAL_ERROR"
)

type APIResponse struct {
	Successful bool             `json:"successful"`
	ErrorCode  string           `json:"error_code"`
	Data       *ProductResponse `json:"data"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	SalePrice   *float64  `json:"sale_price"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toProductResponse(p domain.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		SalePrice:   p.SalePrice,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type PatchAPIResponse struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeCreate(w http.ResponseWriter, status int, successful bool, code string, data *ProductResponse) {
	writeJSON(w, status, APIResponse{
		Successful: successful,
		ErrorCode:  code,
		Data:       data,
	})
}

func writePatch(w http.ResponseWriter, status int, successful bool, code string) {
	writeJSON(w, status, PatchAPIResponse{
		Successful: successful,
		ErrorCode:  code,
	})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest, ErrorValidation
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, ErrorNotFound
	default:
		return http.StatusInternalServerError, ErrorInternal
	}
}
