package catalog

import (
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/utils"
)

type CatalogHandler struct {
	service *CatalogService
}

func NewCatalogHandler(service *CatalogService) *CatalogHandler {
	return &CatalogHandler{
		service: service,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset := utils.ParseIntParam(r.URL.Query().Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}

	limit := utils.ParseIntParam(r.URL.Query().Get("limit"), 10)
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	category := r.URL.Query().Get("category")

	var priceLessThan *float64
	if priceLessThanStr := r.URL.Query().Get("priceLessThan"); priceLessThanStr != "" {
		var price float64
		if _, err := fmt.Sscanf(priceLessThanStr, "%f", &price); err == nil && price > 0 {
			priceLessThan = &price
		}
	}

	response, err := h.service.ListProducts(offset, limit, category, priceLessThan)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.SuccessResponse(w, response)
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	response, err := h.service.GetProductDetails(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	api.SuccessResponse(w, response)
}
