package middleware

import (
	"log"
	"net/http"
)

// Recovery is middleware that recovers from panics and writes a 500 response.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log stack and recovered value, but ensure we don't cause further panics while writing response.
				log.Printf("panic recovered: %v", rec)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal server error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
