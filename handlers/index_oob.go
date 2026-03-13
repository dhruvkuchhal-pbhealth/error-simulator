package handlers

import (
	"net/http"
	"sort"

	"github.com/your-org/error-simulator/models"
)

// ReportGenerator builds reports from product data.
// The bug: GetTopProducts accesses hardcoded indices [0], [1], [5] without bounds checking.
type ReportGenerator struct {
	reportID string
}

// NewReportGenerator returns a report generator.
func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{reportID: "rpt-daily"}
}

// GetTopProducts returns the top N products by revenue. It assumes the slice
// has at least 6 elements and accesses products[0], products[1], products[5]
// without checking length — index [5] panics when len(products) < 6.
func (g *ReportGenerator) GetTopProducts(products []models.Product, n int) []models.Product {
	if n <= 0 {
		return nil
	}
	// Simulate sorting by price * stock (revenue potential).
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price*float64(products[i].Stock) > products[j].Price*float64(products[j].Stock)
	})
	// BUG: Hardcoded indices with no bounds check. If products has fewer than 6
	// elements, products[5] causes index out of range.
	first := products[0]
	second := products[1]
	fifth := products[5]
	return []models.Product{first, second, fifth}
}

// IndexOOB handles GET /error/index-oob.
// It passes a slice of length 3 and calls GetTopProducts which accesses index 5.
func IndexOOB(gen *ReportGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products := []models.Product{
			{ID: "p1", Name: "Widget A", Price: 10.0, Stock: 100},
			{ID: "p2", Name: "Widget B", Price: 20.0, Stock: 50},
			{ID: "p3", Name: "Widget C", Price: 5.0, Stock: 200},
		}
		_ = gen.GetTopProducts(products, 3)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
