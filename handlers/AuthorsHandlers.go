package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
	"um6p.ma/finalproject/validation"
)

type AuthorHandler struct {
	Store interfaces.AuthorStore
}

// GetAuthorByIDHandler retrieves an author by ID from the database
func (h *AuthorHandler) GetAuthorByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid ID format",
		))
		return
	}

	var author models.Author
	if err := database.DB.WithContext(ctx).First(&author, id).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewNotFoundError("Author", id))
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
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request body",
		).WithDebug(err.Error()))
		return
	}

	// Validate author data
	if errors := validation.Validate(newAuthor); len(errors) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errors))
		return
	}

	// Insert into the database
	if err := database.DB.WithContext(ctx).Create(&newAuthor).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
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
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid ID format",
		))
		return
	}

	var updatedAuthor models.Author
	if err := json.NewDecoder(r.Body).Decode(&updatedAuthor); err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request body",
		).WithDebug(err.Error()))
		return
	}

	// Validate author data
	if errors := validation.Validate(updatedAuthor); len(errors) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errors))
		return
	}

	// Check if author exists
	var existingAuthor models.Author
	if err := database.DB.WithContext(ctx).First(&existingAuthor, id).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewNotFoundError("Author", id))
		return
	}

	// Update in the database
	if err := database.DB.WithContext(ctx).Model(&models.Author{}).Where("id = ?", id).Updates(updatedAuthor).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
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
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid ID format",
		))
		return
	}

	// Check if author exists and has no associated books
	var bookCount int64
	if err := database.DB.WithContext(ctx).Model(&models.Book{}).Where("author_id = ?", id).Count(&bookCount).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	if bookCount > 0 {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusConflict,
			errorhandling.ErrCodeBadRequest,
			"Cannot delete author with existing books",
		).WithDetails(map[string]interface{}{
			"bookCount": bookCount,
		}))
		return
	}

	// Delete from the database
	if err := database.DB.WithContext(ctx).Delete(&models.Author{}, id).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListAuthorsHandler retrieves all authors from the database
func (h *AuthorHandler) ListAuthorsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var authors []models.Author

	if err := database.DB.WithContext(ctx).Find(&authors).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}
