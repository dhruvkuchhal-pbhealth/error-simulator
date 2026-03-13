package handlers

import (
	"log"
	"net/http"
	"time"
)

// MetricsService computes business metrics. The bug: CalculateConversionRate
// divides by totalVisits without checking for zero.
type MetricsService struct {
	period string
}

// NewMetricsService returns a metrics service for the given period.
func NewMetricsService(period string) *MetricsService {
	return &MetricsService{period: period}
}

// CalculateConversionRate returns conversions per visit as a rate.
// When totalVisits is 0, conversions/totalVisits causes integer divide by zero.
func (m *MetricsService) CalculateConversionRate(totalVisits int, conversions int) float64 {
	log.Printf("[MetricsService] computing conversion rate for period=%s visits=%d conversions=%d",
		m.period, totalVisits, conversions)
	// BUG: No zero check on totalVisits.
	rate := conversions / totalVisits
	return float64(rate)
}

// DivisionZero handles GET /error/division-zero.
// It calls CalculateConversionRate with totalVisits=0 to trigger divide by zero.
func DivisionZero(svc *MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		period := time.Now().Format("2006-01")
		ms := NewMetricsService(period)
		_ = ms.CalculateConversionRate(0, 5)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
