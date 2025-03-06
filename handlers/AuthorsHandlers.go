package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

type AuthorHandler struct {
	Store interfaces.AuthorStore
}

// GetAuthorByIDHandler retrieves an author by ID from the database
func (h *AuthorHandler) GetAuthorByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var author models.Author
	if err := database.DB.WithContext(ctx).First(&author, id).Error; err != nil {
		http.Error(w, "Author not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

// CreateAuthorHandler adds a new author to the database
func (h *AuthorHandler) CreateAuthorHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newAuthor models.Author
	if err := json.NewDecoder(r.Body).Decode(&newAuthor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert into the database
	if err := database.DB.WithContext(ctx).Create(&newAuthor).Error; err != nil {
		http.Error(w, "Failed to add author", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAuthor)
}

// UpdateAuthorHandler modifies an existing author in the database
func (h *AuthorHandler) UpdateAuthorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedAuthor models.Author
	if err := json.NewDecoder(r.Body).Decode(&updatedAuthor); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update in the database
	if err := database.DB.WithContext(ctx).Model(&models.Author{}).Where("id = ?", id).Updates(updatedAuthor).Error; err != nil {
		http.Error(w, "Author not found or update failed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAuthor)
}

// DeleteAuthorHandler removes an author from the database
func (h *AuthorHandler) DeleteAuthorHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete from the database
	if err := database.DB.WithContext(ctx).Delete(&models.Author{}, id).Error; err != nil {
		http.Error(w, "Author not found or delete failed", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListAuthorsHandler retrieves all authors from the database
func (h *AuthorHandler) ListAuthorsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var authors []models.Author

	if err := database.DB.WithContext(ctx).Find(&authors).Error; err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}
