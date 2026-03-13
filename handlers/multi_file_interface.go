package handlers

import (
	"net/http"

	"github.com/your-org/error-simulator/usersvc"
)

// MultiFileInterface handles GET /error/multi-file/interface.
// Genre: interface boundary — handler calls usersvc.GetUser → impl (userfetcher) panics in its FetchUser.
func MultiFileInterface(svc *usersvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = svc.GetUser("user-interface-1")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
