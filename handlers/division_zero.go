package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

// MetricsService provides methods for metrics calculations.
type MetricsService struct{}

// CalculateConversionRate calculates conversion rate given numerator and denominator.
// Returns an error if the denominator is zero.
func (m *MetricsService) CalculateConversionRate(numerator, denominator int) (int, error) {
	if denominator == 0 {
		return 0, errors.New("denominator is zero")
	}
	return (numerator * 100) / denominator, nil
}

// DivisionZero handler demonstrates a division operation via MetricsService.
func DivisionZero(svc *MetricsService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For the example we use fixed values; in real code these would come from request.
		numerator := 0
		denominator := 5

		// Call the service and handle potential division-by-zero error.
		rate, err := svc.CalculateConversionRate(numerator, denominator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]int{"conversion_rate": rate})
	})
}
