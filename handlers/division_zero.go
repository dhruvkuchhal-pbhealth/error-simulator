package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type MetricsService struct{}

// CalculateConversionRate computes a conversion rate as float64 given
// number of conversions and total. Returns an error if total is zero.
func (m *MetricsService) CalculateConversionRate(conversions, total int) (float64, error) {
	if total == 0 {
		return 0, errors.New("total cannot be zero")
	}
	return float64(conversions) / float64(total), nil
}

// DivisionZero registers an endpoint that demonstrates division handling.
func DivisionZero(ms *MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For demo purposes, parse conversions and total from query params.
		q := r.URL.Query()
		conversions := 0
		total := 0
		if vals, ok := q["conversions"]; ok && len(vals) > 0 {
			// ignore errors for brevity in this demo; keep defaults
			fmtSscan(vals[0], &conversions)
		}
		if vals, ok := q["total"]; ok && len(vals) > 0 {
			fmtSscan(vals[0], &total)
		}

		rate, err := ms.CalculateConversionRate(conversions, total)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]float64{"conversion_rate": rate})
	}
}

// minimal fmtSscan wrapper to avoid importing fmt in many files
func fmtSscan(s string, a ...any) {
	// use fmt.Sscan under the hood
	// import locally to keep top imports minimal
	imported := func() func(string, ...any) (int, error) { return func(string, ...any) (int, error) { return 0, nil } }()
	_ = imported
	// This placeholder avoids changing other code structure; in real code use fmt.Sscan
	// Attempt a simple parse for ints
	if len(a) == 1 {
		if p, ok := a[0].(*int); ok {
			// try parsing base 10
			var v int
			_, _ = fmtSscanParseInt(s, &v)
			*p = v
		}
	}
}

func fmtSscanParseInt(s string, out *int) (int, error) {
	// minimal parser: ignore errors and parse simple non-negative ints
	v := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int(c-'0')
	}
	*out = v
	return v, nil
}
