package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

type AuthorHandler struct {
	Store interfaces.AuthorStore
}

func (h *AuthorHandler) GetAuthorByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	author, err := h.Store.GetAuthor(ctx, id)
	if err != nil {
		http.Error(w, "Author not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func (h *AuthorHandler) CreateAuthorHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newAuthor models.Author
	err := json.NewDecoder(r.Body).Decode(&newAuthor)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdAuthor, err := h.Store.CreateAuthor(ctx, newAuthor)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAuthor)
}

func (h *AuthorHandler) UpdateAuthorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedAuthor models.Author
	err = json.NewDecoder(r.Body).Decode(&updatedAuthor)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	author, err := h.Store.UpdateAuthor(ctx, id, updatedAuthor)
	if err != nil {
		http.Error(w, "Author not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func (h *AuthorHandler) DeleteAuthorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Store.DeleteAuthor(ctx, id)
	if err != nil {
		http.Error(w, "Author not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthorHandler) ListAuthorsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	authors, err := h.Store.GetAllAuthors(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}
