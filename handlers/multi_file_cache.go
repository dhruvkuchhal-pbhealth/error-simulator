package handlers

import (
	"net/http"

	"github.com/your-org/error-simulator/cachesvc"
)

// MultiFileCache handles GET /error/multi-file/cache.
// Call chain: handler → cachesvc.GetUserByID → repo.FindByID (panic in repo).
func MultiFileCache(svc *cachesvc.CacheService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = svc.GetUserByID("user-mf-1")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
