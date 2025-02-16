package inmemorystores

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"

	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/models"
)

type InMemoryBookStore struct {
	mu     sync.RWMutex
	books  map[int]models.Book
	nextID int
}

func NewInMemoryBookStore() *InMemoryBookStore {
	return &InMemoryBookStore{
		books:  make(map[int]models.Book),
		nextID: 1,
	}
}

func (store *InMemoryBookStore) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var books []models.Book
		for _, book := range store.books {
			books = append(books, book)
		}
		return books, nil
	}
}

func (store *InMemoryBookStore) CreateBook(ctx context.Context, book models.Book) (models.Book, error) {
	store.mu.Lock()

	book.ID = store.nextID
	store.nextID++
	store.books[book.ID] = book

	store.mu.Unlock()

	if err := store.SaveBooksToJSON("./database/books.json"); err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (store *InMemoryBookStore) GetBook(ctx context.Context, id int) (models.Book, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return models.Book{}, ctx.Err()
	default:
		book, exists := store.books[id]
		if !exists {
			return models.Book{}, errorhandling.ErrBookNotFound
		}
		return book, nil
	}
}

func (store *InMemoryBookStore) UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error) {
	store.mu.Lock()

	_, exists := store.books[id]
	if !exists {
		return models.Book{}, errorhandling.ErrBookNotFound
	}

	book.ID = id
	store.books[id] = book

	store.mu.Unlock()

	if err := store.SaveBooksToJSON("./database/books.json"); err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (store *InMemoryBookStore) DeleteBook(ctx context.Context, id int) error {
	store.mu.Lock()

	_, exists := store.books[id]
	if !exists {
		return errorhandling.ErrBookNotFound
	}

	delete(store.books, id)

	store.mu.Unlock()

	if err := store.SaveBooksToJSON("./database/books.json"); err != nil {
		return err
	}

	return nil
}

func (store *InMemoryBookStore) SearchBooks(ctx context.Context, criteria models.SearchCriteria) ([]models.Book, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var results []models.Book
	for _, book := range store.books {
		select {
		case <-ctx.Done():
			return nil, ctx.Err() // Handle context cancellation
		default:
		}

		if len(criteria.Titles) > 0 {
			for _, title := range criteria.Titles {
				if book.Title == title {
					results = append(results, book)
				}
			}
		} else if len(criteria.Authors) > 0 {
			for _, author := range criteria.Authors {
				if book.Author.FirstName == author || book.Author.LastName == author {
					results = append(results, book)
				}
			}
		} else if len(criteria.Genres) > 0 {
			for _, genre := range criteria.Genres {
				bookGenres := strings.Split(book.Genres, ", ")

				for _, bookGenre := range bookGenres {
					if genre == bookGenre {
						results = append(results, book)
						break
					}
				}
			}
		} else if criteria.MinPrice >= 0 && criteria.MaxPrice > 0 {
			if book.Price >= criteria.MinPrice && book.Price <= criteria.MaxPrice {
				results = append(results, book)
			}
		}
	}

	return results, nil
}

func (store *InMemoryBookStore) LoadBooksFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var books []models.Book
	if err := json.NewDecoder(file).Decode(&books); err != nil {
		return err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for _, book := range books {
		store.books[book.ID] = book
		if book.ID >= store.nextID {
			store.nextID = book.ID + 1
		}
	}

	return nil
}

func (store *InMemoryBookStore) SaveBooksToJSON(filePath string) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var books []models.Book
	for _, book := range store.books {
		books = append(books, book)
	}

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
