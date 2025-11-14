package category

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type CategoriesHandler struct {
	service *CategoriesService
}

func NewCategoriesHandler(service *CategoriesService) *CategoriesHandler {
	return &CategoriesHandler{
		service: service,
	}
}

func (h *CategoriesHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.ListCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.SuccessResponse(w, response)
}

func (h *CategoriesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.CreateCategory(req)
	if err != nil {
		if errors.Is(err, ErrCategoryCodeRequired) ||
			errors.Is(err, ErrCategoryNameRequired) ||
			errors.Is(err, ErrCategoryCodeTooLong) ||
			errors.Is(err, ErrCategoryNameTooLong) {
			api.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.CreatedResponse(w, response)
}
