package inmemorystores

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"um6p.ma/finalproject/errorhandling"
	"um6p.ma/finalproject/models"
)

type InMemoryCustomerStore struct {
	mu        sync.RWMutex
	customers map[int]models.Customer
	nextID    int
}

func NewInMemoryCustomerStore() *InMemoryCustomerStore {
	return &InMemoryCustomerStore{
		customers: make(map[int]models.Customer),
		nextID:    1,
	}
}

func (store *InMemoryCustomerStore) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var customers []models.Customer
		for _, customer := range store.customers {
			customers = append(customers, customer)
		}
		return customers, nil
	}
}

func (store *InMemoryCustomerStore) CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return models.Customer{}, ctx.Err()
	default:
		customer.ID = store.nextID
		store.nextID++
		store.customers[customer.ID] = customer

		store.mu.Unlock()

		if err := store.SaveCustomersToJSON("./database/customers.json"); err != nil {
			return models.Customer{}, err
		}

		return customer, nil
	}
}

func (store *InMemoryCustomerStore) GetCustomer(ctx context.Context, id int) (models.Customer, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	select {
	case <-ctx.Done():
		return models.Customer{}, ctx.Err()
	default:
		customer, exists := store.customers[id]
		if !exists {
			return models.Customer{}, errorhandling.ErrCustomerNotFound
		}
		return customer, nil
	}
}

func (store *InMemoryCustomerStore) UpdateCustomer(ctx context.Context, id int, customer models.Customer) (models.Customer, error) {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return models.Customer{}, ctx.Err()
	default:
		_, exists := store.customers[id]
		if !exists {
			return models.Customer{}, errorhandling.ErrCustomerNotFound
		}

		customer.ID = id
		store.customers[id] = customer

		store.mu.Unlock()

		if err := store.SaveCustomersToJSON("./database/customers.json"); err != nil {
			return models.Customer{}, err
		}

		return customer, nil
	}
}

func (store *InMemoryCustomerStore) DeleteCustomer(ctx context.Context, id int) error {
	store.mu.Lock()

	select {
	case <-ctx.Done():
		store.mu.Unlock()
		return ctx.Err()
	default:
		_, exists := store.customers[id]
		if !exists {
			return errorhandling.ErrCustomerNotFound
		}

		delete(store.customers, id)
		store.mu.Unlock()

		if err := store.SaveCustomersToJSON("./database/customers.json"); err != nil {
			return err
		}

		return nil
	}
}

func (store *InMemoryCustomerStore) LoadCustomersFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var customers []models.Customer
	if err := json.NewDecoder(file).Decode(&customers); err != nil {
		return err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	for _, customer := range customers {
		store.customers[customer.ID] = customer
		if customer.ID >= store.nextID {
			store.nextID = customer.ID + 1
		}
	}

	return nil
}

func (store *InMemoryCustomerStore) SaveCustomersToJSON(filePath string) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var customers []models.Customer
	for _, customer := range store.customers {
		customers = append(customers, customer)
	}

	data, err := json.MarshalIndent(customers, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
