package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/Elvi-2202/book-api/internal/model"
	"github.com/Elvi-2202/book-api/internal/repository"
)

type BookHandler struct {
	repo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b model.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		sendError(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	if err := b.Validate(); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	if err := h.repo.Create(r.Context(), &b); err != nil {
		sendError(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) 
	json.NewEncoder(w).Encode(b)
}

func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.List(r.Context())
	if err != nil {
		sendError(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := h.repo.GetByID(r.Context(), id)
	
	if err != nil {
		sendError(w, "Livre non trouvé", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var b model.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		sendError(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	b.ID = id
	if err := b.Validate(); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	b.UpdatedAt = time.Now()

	if err := h.repo.Update(r.Context(), &b); err != nil {
		sendError(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		sendError(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)// Respect du code 204 [cite: 45]
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}