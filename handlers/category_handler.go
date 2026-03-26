package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"Book-API_Golang/models"
)

type CategoryHandler struct {
	Store *models.Store
}

func (h *CategoryHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/categories", h.ListCategories).Methods(http.MethodGet)
	r.HandleFunc("/categories", h.CreateCategory).Methods(http.MethodPost)
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.ListCategories())
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var in nameInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	category, err := h.Store.CreateCategory(in.Name)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrValidation):
			writeError(w, http.StatusBadRequest, "name is required")
		case errors.Is(err, models.ErrDuplicateName):
			writeError(w, http.StatusConflict, "category already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, category)
}
