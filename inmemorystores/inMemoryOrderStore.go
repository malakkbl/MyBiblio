package inmemorystores

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/models"
)

type InMemoryOrderStore struct {
	mu        sync.RWMutex
	orders    map[int]models.Order
	nextID    int
	bookStore *InMemoryBookStore
}

func NewInMemoryOrderStore(bookStore *InMemoryBookStore) *InMemoryOrderStore {
	return &InMemoryOrderStore{
		orders:    make(map[int]models.Order),
		nextID:    1,
		bookStore: bookStore,
	}
}

func (store *InMemoryOrderStore) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return models.Order{}, ctx.Err()
	default:
	}

	for _, item := range order.Items {
		book, err := store.bookStore.GetBook(ctx, item.Book.ID)
		if err != nil {
			store.mu.Unlock()
			return models.Order{}, errors.New("book does not exist")
		}
		if book.Stock < item.Quantity {
			store.mu.Unlock()
			return models.Order{}, errors.New("insufficient stock for book")
		}
		book.Stock -= item.Quantity
		_, err = store.bookStore.UpdateBook(ctx, book.ID, book)
		if err != nil {
			store.mu.Unlock()
			log.Printf("you re here")
			return models.Order{}, err
		}
	}

	order.ID = store.nextID
	store.nextID++
	store.orders[order.ID] = order

	store.mu.Unlock()

	if err := store.SaveOrdersToJSON("./database/orders.json"); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (store *InMemoryOrderStore) GetOrder(ctx context.Context, id int) (models.Order, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
	}

	order, exists := store.orders[id]
	if !exists {
		return models.Order{}, errorhandling.ErrOrderNotFound
	}

	return order, nil
}

func (store *InMemoryOrderStore) UpdateOrder(ctx context.Context, id int, order models.Order) (models.Order, error) {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return models.Order{}, ctx.Err()
	default:
	}

	existingOrder, exists := store.orders[id]
	if !exists {
		store.mu.Unlock()
		return models.Order{}, errorhandling.ErrOrderNotFound
	}

	for _, item := range existingOrder.Items {
		book, err := store.bookStore.GetBook(ctx, item.Book.ID)
		if err != nil {
			store.mu.Unlock()
			return models.Order{}, errors.New("book does not exist during stock restoration")
		}
		book.Stock += item.Quantity
		_, err = store.bookStore.UpdateBook(ctx, book.ID, book)
		if err != nil {
			store.mu.Unlock()
			return models.Order{}, err
		}
	}

	for _, item := range order.Items {
		book, err := store.bookStore.GetBook(ctx, item.Book.ID)
		if err != nil {
			store.mu.Unlock()
			return models.Order{}, errors.New("book does not exist during stock validation")
		}
		if book.Stock < item.Quantity {
			store.mu.Unlock()
			return models.Order{}, errorhandling.ErrInsufficientStock
		}
		book.Stock -= item.Quantity
		_, err = store.bookStore.UpdateBook(ctx, book.ID, book)
		if err != nil {
			store.mu.Unlock()
			return models.Order{}, err
		}
	}

	order.ID = id
	store.orders[id] = order

	store.mu.Unlock()

	if err := store.SaveOrdersToJSON("./database/orders.json"); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (store *InMemoryOrderStore) DeleteOrder(ctx context.Context, id int) error {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return ctx.Err()
	default:
	}

	_, exists := store.orders[id]
	if !exists {
		store.mu.Unlock()
		return errorhandling.ErrOrderNotFound
	}

	delete(store.orders, id)

	store.mu.Unlock()

	if err := store.SaveOrdersToJSON("./database/orders.json"); err != nil {
		return err
	}

	return nil
}

func (store *InMemoryOrderStore) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var orders []models.Order
	for _, order := range store.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

func (store *InMemoryOrderStore) LoadOrdersFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var orders []models.Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for _, order := range orders {
		store.orders[order.ID] = order
		if order.ID >= store.nextID {
			store.nextID = order.ID + 1
		}
	}

	return nil
}

func (store *InMemoryOrderStore) SaveOrdersToJSON(filePath string) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var orders []models.Order
	for _, order := range store.orders {
		orders = append(orders, order)
	}

	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func (store *InMemoryOrderStore) GetOrdersInTimeRange(ctx context.Context, start, end time.Time) ([]models.Order, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var ordersInRange []models.Order
	for _, order := range store.orders {
		if order.CreatedAt.After(start) && order.CreatedAt.Before(end) {
			ordersInRange = append(ordersInRange, order)
		}
	}

	return ordersInRange, nil
}
