package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/handlers"
	"um6p.ma/finalproject/inmemorystores"
	"um6p.ma/finalproject/interfaces"
	"um6p.ma/finalproject/models"
)

var (
	latestSalesReport models.SalesReport
	reportMutex       sync.RWMutex
)

func main() {
	bookStore := inmemorystores.NewInMemoryBookStore()
	authorStore := inmemorystores.NewInMemoryAuthorStore()
	customerStore := inmemorystores.NewInMemoryCustomerStore()
	orderStore := inmemorystores.NewInMemoryOrderStore(bookStore)

	if err := bookStore.LoadBooksFromJSON("./database/books.json"); err != nil {
		log.Printf("Failed to load books: %v", err)
	}
	if err := authorStore.LoadAuthorsFromJSON("./database/authors.json"); err != nil {
		log.Printf("Failed to load authors: %v", err)
	}
	if err := customerStore.LoadCustomersFromJSON("./database/customers.json"); err != nil {
		log.Printf("Failed to load customers: %v", err)
	}
	if err := orderStore.LoadOrdersFromJSON("./database/orders.json"); err != nil {
		log.Printf("Failed to load orders: %v", err)
	}

	bookHandler := &handlers.BookHandler{Store: bookStore}
	authorHandler := &handlers.AuthorHandler{Store: authorStore}
	customerHandler := &handlers.CustomerHandler{Store: customerStore}
	orderHandler := &handlers.OrderHandler{Store: orderStore}

	router := httprouter.New()
	router.GET("/books/:id", bookHandler.GetBookByIDHandler)
	router.POST("/books", bookHandler.CreateBookHandler)
	router.PUT("/books/:id", bookHandler.UpdateBookHandler)
	router.DELETE("/books/:id", bookHandler.DeleteBookHandler)
	router.GET("/books", bookHandler.SearchBooksHandler)
	router.GET("/authors/:id", authorHandler.GetAuthorByIDHandler)
	router.POST("/authors", authorHandler.CreateAuthorHandler)
	router.PUT("/authors/:id", authorHandler.UpdateAuthorHandler)
	router.DELETE("/authors/:id", authorHandler.DeleteAuthorHandler)
	router.GET("/authors", authorHandler.ListAuthorsHandler)
	router.GET("/customers/:id", customerHandler.GetCustomerByIDHandler)
	router.POST("/customers", customerHandler.CreateCustomerHandler)
	router.PUT("/customers/:id", customerHandler.UpdateCustomerHandler)
	router.DELETE("/customers/:id", customerHandler.DeleteCustomerHandler)
	router.GET("/customers", customerHandler.ListCustomersHandler)
	router.GET("/orders", orderHandler.GetAllOrdersHandler)
	router.GET("/orders/:id", orderHandler.GetOrderByIDHandler)
	router.POST("/orders", orderHandler.CreateOrderHandler)
	router.PUT("/orders/:id", orderHandler.UpdateOrderHandler)
	router.DELETE("/orders/:id", orderHandler.DeleteOrderHandler)

	router.GET("/sales-reports", getSalesReportHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		<-stop
		log.Println("Shutting down gracefully...")
		cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error during server shutdown: %v", err)
		}
	}()

	go startSalesReportGeneration(ctx, orderStore, bookStore)

	log.Println("Starting server on port 8080...")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped.")
}

func generateSalesReport(ctx context.Context, orderStore interfaces.OrderStore, bookStore interfaces.BookStore) (models.SalesReport, error) {
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	orders, err := orderStore.GetOrdersInTimeRange(ctx, startTime, endTime)
	if err != nil {
		return models.SalesReport{}, err
	}

	totalRevenue := 0.0
	totalOrders := len(orders)
	bookSalesMap := make(map[int]int)

	for _, order := range orders {
		totalRevenue += order.TotalPrice
		for _, item := range order.Items {
			bookSalesMap[item.Book.ID] += item.Quantity
		}
	}

	topSellingBooks := calculateTopSellingBooks(bookSalesMap, bookStore)

	report := models.SalesReport{
		Timestamp:       time.Now(),
		TotalRevenue:    totalRevenue,
		TotalOrders:     totalOrders,
		TopSellingBooks: topSellingBooks,
	}

	reportMutex.Lock()
	latestSalesReport = report
	reportMutex.Unlock()

	err = saveSalesReportToFile(report)
	if err != nil {
		return models.SalesReport{}, err
	}

	return report, nil
}

func calculateTopSellingBooks(bookSalesMap map[int]int, bookStore interfaces.BookStore) []models.BookSales {
	var bookSales []models.BookSales

	for bookID, quantity := range bookSalesMap {
		book, err := bookStore.GetBook(context.Background(), bookID)
		if err == nil {
			bookSales = append(bookSales, models.BookSales{
				Book:     book,
				Quantity: quantity,
			})
		}
	}

	sort.Slice(bookSales, func(i, j int) bool {
		return bookSales[i].Quantity > bookSales[j].Quantity
	})

	if len(bookSales) > 10 {
		bookSales = bookSales[:10]
	}

	return bookSales
}

func saveSalesReportToFile(report models.SalesReport) error {
	filePath := "./database/sales_reports.json"

	var reports []models.SalesReport

	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&reports); err != nil {
			return err
		}
	}

	reports = append(reports, report)

	data, err := json.MarshalIndent(reports, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func startSalesReportGeneration(ctx context.Context, orderStore interfaces.OrderStore, bookStore interfaces.BookStore) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Sales report generation stopped.")
			return
		case <-ticker.C:
			_, err := generateSalesReport(ctx, orderStore, bookStore)
			if err != nil {
				log.Printf("Error generating sales report: %v", err)
			}
		}
	}
}

func getSalesReportHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filePath := "./database/sales_reports.json"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "No sales reports available", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var reports []models.SalesReport
	if err := json.NewDecoder(file).Decode(&reports); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
