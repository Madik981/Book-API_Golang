package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"Book-API_Golang/models"
)

type BookHandler struct {
	Store *models.Store
}

type listBooksResponse struct {
	Data       []models.Book `json:"data"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	Total      int           `json:"total"`
	TotalPages int           `json:"total_pages"`
}

func (h *BookHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/books", h.ListBooks).Methods(http.MethodGet)
	r.HandleFunc("/books", h.CreateBook).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", h.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", h.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books/{id}", h.DeleteBook).Methods(http.MethodDelete)
}

func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page := parsePositiveInt(q.Get("page"), 1)
	pageSize := parsePositiveInt(q.Get("page_size"), 10)
	authorID := parsePositiveInt(q.Get("author_id"), 0)
	categoryID := parsePositiveInt(q.Get("category_id"), 0)

	books := h.Store.ListBooks(authorID, categoryID)
	total := len(books)
	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	resp := listBooksResponse{
		Data:       books[start:end],
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	book, err := h.Store.GetBook(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeError(w, http.StatusNotFound, "book not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var in models.Book
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	book, err := h.Store.CreateBook(in)
	if err != nil {
		handleBookMutationError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var in models.Book
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	book, err := h.Store.UpdateBook(id, in)
	if err != nil {
		handleBookMutationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err := h.Store.DeleteBook(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeError(w, http.StatusNotFound, "book not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleBookMutationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrValidation):
		writeError(w, http.StatusBadRequest, "validation failed")
	case errors.Is(err, models.ErrRelationAbsent):
		writeError(w, http.StatusBadRequest, "author_id or category_id does not exist")
	case errors.Is(err, models.ErrNotFound):
		writeError(w, http.StatusNotFound, "book not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func pathID(r *http.Request) (int, bool) {
	raw := mux.Vars(r)["id"]
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
