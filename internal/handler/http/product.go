package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/demo/product-api/internal/domain"
)

type ProductService interface {
	Create(ctx context.Context, in domain.CreateProductInput) (domain.Product, error)
	Patch(ctx context.Context, id string, in domain.PatchProductInput) error
}

type ProductHandler struct {
	svc ProductService
}

func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}

type PatchProductRequest struct {
	Name        json.RawMessage `json:"name"`
	Description json.RawMessage `json:"description"`
	SalePrice   json.RawMessage `json:"sale_price"`
	Price       json.RawMessage `json:"price"`
}

// CreateProduct godoc
// @Summary Create product
// @Tags product
// @Accept json
// @Produce json
// @Param body body CreateProductRequest true "Create product payload"
// @Success 201 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /product [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeCreate(w, http.StatusBadRequest, false, ErrorValidation, nil)
		return
	}
	created, err := h.svc.Create(r.Context(), domain.CreateProductInput{
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	})
	if err != nil {
		status, code := mapError(err)
		writeCreate(w, status, false, code, nil)
		return
	}
	writeCreate(w, http.StatusCreated, true, ErrorNone, &CreateData{
		Data1: created.ID,
		Data2: created.Name,
	})
}

// PatchProduct godoc
// @Summary Patch product
// @Description Partial update. Omitted fields are left unchanged. description and sale_price accept null.
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param body body PatchProductRequest true "Patch product payload"
// @Success 200 {object} PatchAPIResponse
// @Failure 400 {object} PatchAPIResponse
// @Failure 404 {object} PatchAPIResponse
// @Failure 500 {object} PatchAPIResponse
// @Router /product/{id} [patch]
func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req PatchProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writePatch(w, http.StatusBadRequest, false, ErrorValidation)
		return
	}
	in, err := req.toInput()
	if err != nil {
		writePatch(w, http.StatusBadRequest, false, ErrorValidation)
		return
	}
	if err := h.svc.Patch(r.Context(), id, in); err != nil {
		status, code := mapError(err)
		writePatch(w, status, false, code)
		return
	}
	writePatch(w, http.StatusOK, true, ErrorNone)
}

func (req PatchProductRequest) toInput() (domain.PatchProductInput, error) {
	name, err := decodeOptionalString(req.Name, false)
	if err != nil {
		return domain.PatchProductInput{}, err
	}
	desc, err := decodeOptionalString(req.Description, true)
	if err != nil {
		return domain.PatchProductInput{}, err
	}
	sale, err := decodeOptionalFloat(req.SalePrice, true)
	if err != nil {
		return domain.PatchProductInput{}, err
	}
	price, err := decodeOptionalFloat(req.Price, false)
	if err != nil {
		return domain.PatchProductInput{}, err
	}
	return domain.PatchProductInput{
		Name:        name,
		Description: desc,
		SalePrice:   sale,
		Price:       price,
	}, nil
}

func decodeOptionalString(raw json.RawMessage, nullable bool) (domain.Optional[string], error) {
	if len(raw) == 0 {
		return domain.Optional[string]{}, nil
	}
	if string(raw) == "null" {
		if !nullable {
			return domain.Optional[string]{}, fmt.Errorf("field cannot be null")
		}
		return domain.Optional[string]{Set: true, Value: nil}, nil
	}
	var v string
	if err := json.Unmarshal(raw, &v); err != nil {
		return domain.Optional[string]{}, err
	}
	return domain.Optional[string]{Set: true, Value: &v}, nil
}

func decodeOptionalFloat(raw json.RawMessage, nullable bool) (domain.Optional[float64], error) {
	if len(raw) == 0 {
		return domain.Optional[float64]{}, nil
	}
	if string(raw) == "null" {
		if !nullable {
			return domain.Optional[float64]{}, fmt.Errorf("field cannot be null")
		}
		return domain.Optional[float64]{Set: true, Value: nil}, nil
	}
	var v float64
	if err := json.Unmarshal(raw, &v); err != nil {
		return domain.Optional[float64]{}, err
	}
	return domain.Optional[float64]{Set: true, Value: &v}, nil
}

func decodeJSON(r *http.Request, dest any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return err
	}
	return nil
}
