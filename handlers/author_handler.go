package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"Book-API_Golang/models"
)

type AuthorHandler struct {
	Store *models.Store
}

type nameInput struct {
	Name string `json:"name"`
}

func (h *AuthorHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/authors", h.ListAuthors).Methods(http.MethodGet)
	r.HandleFunc("/authors", h.CreateAuthor).Methods(http.MethodPost)
}

func (h *AuthorHandler) ListAuthors(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.ListAuthors())
}

func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var in nameInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	author, err := h.Store.CreateAuthor(in.Name)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrValidation):
			writeError(w, http.StatusBadRequest, "name is required")
		case errors.Is(err, models.ErrDuplicateName):
			writeError(w, http.StatusConflict, "author already exists")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, author)
}
