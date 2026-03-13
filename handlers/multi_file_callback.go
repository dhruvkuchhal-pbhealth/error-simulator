package handlers

import (
	"net/http"

	"github.com/your-org/error-simulator/processor"
)

// MultiFileCallback handles GET /error/multi-file/callback.
// Genre: callback/visitor — we pass a callback to processor.Process; our callback panics (nil Child).
func MultiFileCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := []processor.Item{
			{Name: "a", Child: nil}, // callback will deref Child → panic in this handler's closure
		}
		processor.Process(items, func(it processor.Item) {
			_ = it.Child.Name // panic: it.Child is nil (genre: callback — panic in caller's code)
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
