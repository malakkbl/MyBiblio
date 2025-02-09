package inmemorystores

import (
	"sync"

	"um6p.ma/finalproject/models"
)

type ReportStore struct {
	mu     sync.RWMutex
	report models.SalesReport
}

func (rs *ReportStore) StoreReport(report models.SalesReport) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.report = report
}

func (rs *ReportStore) GetLatestReport() models.SalesReport {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.report
}
