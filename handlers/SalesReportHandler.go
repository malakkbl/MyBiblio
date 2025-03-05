package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/julienschmidt/httprouter"
	"um6p.ma/finalproject/database"
	"um6p.ma/finalproject/models"
)

// GenerateSalesReport generates a report based on database data
func GenerateSalesReport(ctx context.Context) (models.SalesReport, error) {
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	// Fetch all orders from the last 24 hours
	var orders []models.Order
	if err := database.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&orders).Error; err != nil {
		return models.SalesReport{}, err
	}

	totalRevenue := 0.0
	totalOrders := len(orders)
	bookSalesMap := make(map[int]int)

	// Calculate total revenue and book sales
	for _, order := range orders {
		totalRevenue += order.TotalPrice
		for _, item := range order.Items {
			bookSalesMap[item.BookID] += item.Quantity
		}
	}

	// Get top-selling books
	topSellingBooks, err := calculateTopSellingBooks(bookSalesMap)
	if err != nil {
		return models.SalesReport{}, err
	}

	// Generate report
	report := models.SalesReport{
		Timestamp:       time.Now(),
		TotalRevenue:    totalRevenue,
		TotalOrders:     totalOrders,
		TopSellingBooks: topSellingBooks,
	}

	// Store report in the database
	if err := database.DB.Create(&report).Error; err != nil {
		return models.SalesReport{}, err
	}

	log.Println("Sales report generated successfully!")
	return report, nil
}

// calculateTopSellingBooks fetches book details for top-selling books
func calculateTopSellingBooks(bookSalesMap map[int]int) ([]models.BookSales, error) {
	var bookSales []models.BookSales

	for bookID, quantity := range bookSalesMap {
		var book models.Book
		if err := database.DB.First(&book, bookID).Error; err != nil {
			continue
		}
		bookSales = append(bookSales, models.BookSales{
			BookID:   bookID,
			Book:     book,
			Quantity: quantity,
		})
	}

	// Sort books by quantity sold
	sort.Slice(bookSales, func(i, j int) bool {
		return bookSales[i].Quantity > bookSales[j].Quantity
	})

	// Return top 10 books
	if len(bookSales) > 10 {
		bookSales = bookSales[:10]
	}

	return bookSales, nil
}

// GetSalesReportHandler returns the latest reports
func GetSalesReportHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reports []models.SalesReport
	if err := database.DB.Order("timestamp DESC").Limit(10).Find(&reports).Error; err != nil {
		http.Error(w, "Failed to fetch sales reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// StartSalesReportGeneration runs report generation every 24 hours
func StartSalesReportGeneration(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Sales report generation stopped.")
			return
		case <-ticker.C:
			_, err := GenerateSalesReport(ctx)
			if err != nil {
				log.Printf("Error generating sales report: %v", err)
			}
		}
	}
}
