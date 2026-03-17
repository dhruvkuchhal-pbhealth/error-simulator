package handlers

import (
	"encoding/json"
	"net/http"
)

// OrderService simulates an order processing service.
// It may have dependencies in a real application (e.g., DB client).
// For this example, we keep it minimal.
type OrderService struct {
	// placeholder dependency to illustrate nil checks
	Client interface{}
}

// ProcessOrder processes an order. It was panicking when called on a nil
// receiver or when required dependencies were nil. This implementation
// guards against nil receiver and nil dependency and returns an HTTP 500
// with a JSON error instead of causing a panic.
func (svc *OrderService) ProcessOrder(w http.ResponseWriter, r *http.Request) {
	// guard: ensure receiver is non-nil
	if svc == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// guard: ensure required dependency is present
	if svc.Client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "service not initialized"})
		return
	}

	// Normal processing (simulated)
	resp := map[string]string{"status": "processed"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// NilPointer demonstrates an HTTP handler that may call ProcessOrder with a nil service.
// It is left as-is; callers should ensure they pass a properly initialized service.
var NilPointer = struct{
	func1 http.HandlerFunc
}{}

// For completeness, provide a handler function that uses OrderService. In real usage
// the server setup would assign a properly initialized OrderService. We provide a
// safe wrapper here to avoid accidental nil dereference when used directly.

// The original code had a function literal; we provide an exported helper to
// create a safe handler given a possibly-nil *OrderService.
func NewProcessOrderHandler(svc *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If svc is nil, still call ProcessOrder on nil receiver which is now safe
		// because ProcessOrder checks for nil receiver.
		svc.ProcessOrder(w, r)
	}
}
