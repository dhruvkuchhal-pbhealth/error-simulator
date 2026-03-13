package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID retrieves a user by ID and writes the response.
func (r *UserRepository) GetUserByID(w http.ResponseWriter, req *http.Request) {
	// Guard against nil DB to prevent nil pointer dereference panics.
	if r == nil || r.db == nil {
		http.Error(w, "database not initialized", http.StatusInternalServerError)
		return
	}

	idStr := req.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	var name string
	row := r.db.QueryRow("SELECT name FROM users WHERE id = ?", id)
	if err := row.Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"id": id, "name": name}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
