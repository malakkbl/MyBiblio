package inmemorystores

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/models"
)

type InMemoryAuthorStore struct {
	mu      sync.RWMutex
	authors map[int]models.Author
	nextID  int
}

func NewInMemoryAuthorStore() *InMemoryAuthorStore {
	return &InMemoryAuthorStore{
		authors: make(map[int]models.Author),
		nextID:  1,
	}
}

func (store *InMemoryAuthorStore) CreateAuthor(ctx context.Context, author models.Author) (models.Author, error) {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return models.Author{}, ctx.Err()
	default:
	}

	author.ID = store.nextID
	store.nextID++
	store.authors[author.ID] = author

	store.mu.Unlock()

	if err := store.SaveAuthorsToJSON("./database/authors.json"); err != nil {
		return models.Author{}, err
	}

	return author, nil
}

func (store *InMemoryAuthorStore) GetAuthor(ctx context.Context, id int) (models.Author, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}

	author, exists := store.authors[id]
	if !exists {
		return models.Author{}, errorhandling.ErrAuthorNotFound
	}

	return author, nil
}

func (store *InMemoryAuthorStore) UpdateAuthor(ctx context.Context, id int, author models.Author) (models.Author, error) {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return models.Author{}, ctx.Err()
	default:
	}

	_, exists := store.authors[id]
	if !exists {
		return models.Author{}, errorhandling.ErrAuthorNotFound
	}

	author.ID = id
	store.authors[id] = author

	store.mu.Unlock()

	if err := store.SaveAuthorsToJSON("./database/authors.json"); err != nil {
		return models.Author{}, err
	}

	return author, nil
}

func (store *InMemoryAuthorStore) DeleteAuthor(ctx context.Context, id int) error {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return ctx.Err()
	default:
	}

	_, exists := store.authors[id]
	if !exists {
		return errorhandling.ErrAuthorNotFound
	}

	delete(store.authors, id)

	store.mu.Unlock()

	if err := store.SaveAuthorsToJSON("./database/authors.json"); err != nil {
		return err
	}

	return nil
}

func (store *InMemoryAuthorStore) GetAllAuthors(ctx context.Context) ([]models.Author, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var authors []models.Author
	for _, author := range store.authors {
		authors = append(authors, author)
	}

	return authors, nil
}

func (store *InMemoryAuthorStore) LoadAuthorsFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var authors []models.Author
	if err := json.NewDecoder(file).Decode(&authors); err != nil {
		return err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for _, author := range authors {
		store.authors[author.ID] = author
		if author.ID >= store.nextID {
			store.nextID = author.ID + 1
		}
	}

	return nil
}

func (store *InMemoryAuthorStore) SaveAuthorsToJSON(filePath string) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var authors []models.Author
	for _, author := range store.authors {
		authors = append(authors, author)
	}

	data, err := json.MarshalIndent(authors, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
