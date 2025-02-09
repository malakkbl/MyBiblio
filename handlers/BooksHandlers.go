package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

type BookHandler struct {
	Store interfaces.BookStore
}

func (h *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	book, err := h.Store.GetBook(ctx, id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdBook, err := h.Store.CreateBook(ctx, newBook)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBook)
}

func (h *BookHandler) UpdateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedBook models.Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.Store.UpdateBook(ctx, id, updatedBook)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.Store.DeleteBook(ctx, id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) SearchBooksHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	query := r.URL.Query()

	criteria := models.SearchCriteria{
		Titles:  query["title"],
		Authors: query["author"],
		Genres:  query["genre"],
	}

	if minPriceStr := query.Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			criteria.MinPrice = minPrice
		} else {
			http.Error(w, "Invalid min_price value", http.StatusBadRequest)
			return
		}
	}

	if maxPriceStr := query.Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			criteria.MaxPrice = maxPrice
		} else {
			http.Error(w, "Invalid max_price value", http.StatusBadRequest)
			return
		}
	}

	if len(query) == 0 || (len(criteria.Titles) == 0 && len(criteria.Authors) == 0 && len(criteria.Genres) == 0 && criteria.MinPrice == 0 && criteria.MaxPrice == 0) {
		results, err := h.Store.GetAllBooks(ctx)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	results, err := h.Store.SearchBooks(ctx, criteria)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
