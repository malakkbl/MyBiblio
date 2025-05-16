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

type BookHandler struct {
	Store interfaces.BookStore
}

// GetBookByIDHandler retrieves a book by ID from the database
func (h *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var book models.Book
	if err := database.DB.WithContext(ctx).First(&book, id).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewNotFoundError("Book", id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// CreateBookHandler adds a new book to the database
func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newBook models.Book

	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request body",
		).WithDebug(err.Error()))
		return
	}

	// Validate the book data
	if errors := validation.Validate(newBook); len(errors) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errors))
		return
	}

	// Insert into the database
	if err := database.DB.WithContext(ctx).Create(&newBook).Error; err != nil {
		if errorhandling.IsDuplicateKeyError(err) {
			errorhandling.HandleError(w, errorhandling.NewError(
				http.StatusConflict,
				errorhandling.ErrCodeDuplicateEntry,
				"Book already exists",
			))
			return
		}
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// UpdateBookHandler modifies an existing book in the database
func (h *BookHandler) UpdateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		errorhandling.HandleError(w, errorhandling.NewError(
			http.StatusBadRequest,
			errorhandling.ErrCodeInvalidInput,
			"Invalid request body",
		).WithDebug(err.Error()))
		return
	}

	// Validate the book data
	if errors := validation.Validate(updatedBook); len(errors) > 0 {
		errorhandling.HandleError(w, errorhandling.NewValidationError(errors))
		return
	}

	// Update in the database
	if err := database.DB.WithContext(ctx).Model(&models.Book{}).Where("id = ?", id).Updates(updatedBook).Error; err != nil {
		if err := database.DB.WithContext(ctx).First(&models.Book{}, id).Error; err != nil {
			errorhandling.HandleError(w, errorhandling.NewNotFoundError("Book", id))
			return
		}
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

// DeleteBookHandler removes a book from the database
func (h *BookHandler) DeleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	if err := database.DB.WithContext(ctx).Delete(&models.Book{}, id).Error; err != nil {
		if err := database.DB.WithContext(ctx).First(&models.Book{}, id).Error; err != nil {
			errorhandling.HandleError(w, errorhandling.NewNotFoundError("Book", id))
			return
		}
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchBooksHandler searches for books in the database
func (h *BookHandler) SearchBooksHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	query := r.URL.Query()

	var books []models.Book
	dbQuery := database.DB.WithContext(ctx)

	// Apply filters if provided
	if title := query.Get("title"); title != "" {
		dbQuery = dbQuery.Where("title ILIKE ?", "%"+title+"%")
	}
	if authorID := query.Get("author_id"); authorID != "" {
		dbQuery = dbQuery.Where("author_id = ?", authorID)
	}
	if genre := query.Get("genre"); genre != "" {
		dbQuery = dbQuery.Where("genres ILIKE ?", "%"+genre+"%")
	}
	if minPrice := query.Get("min_price"); minPrice != "" {
		dbQuery = dbQuery.Where("price >= ?", minPrice)
	}
	if maxPrice := query.Get("max_price"); maxPrice != "" {
		dbQuery = dbQuery.Where("price <= ?", maxPrice)
	}

	// Execute query
	if err := dbQuery.Find(&books).Error; err != nil {
		errorhandling.HandleError(w, errorhandling.NewDatabaseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
