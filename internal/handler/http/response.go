package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/demo/product-api/internal/domain"
)

const (
	ErrorNone       = ""
	ErrorValidation = "VALIDATION_ERROR"
	ErrorNotFound   = "NOT_FOUND"
	ErrorInternal   = "INTERNAL_ERROR"
)

type APIResponse struct {
	Successful bool        `json:"successful"`
	ErrorCode  string      `json:"error_code"`
	Data       *CreateData `json:"data"`
}

type CreateData struct {
	Data1 string `json:"data1"`
	Data2 string `json:"data2"`
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

func writeCreate(w http.ResponseWriter, status int, successful bool, code string, data *CreateData) {
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
