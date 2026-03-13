package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery middleware recovers from panics in handlers and logs stack traces.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log the panic value and the stack trace. Avoid further panics here.
				log.Printf("panic recovered: %v\n%s", rec, debug.Stack())
				// Return a 500 response if headers not written yet.
				if rw, ok := w.(interface{ WriteHeader(int) }); ok {
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
