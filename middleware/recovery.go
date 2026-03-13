package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/your-org/error-simulator/config"
	"github.com/your-org/error-simulator/kafka"
	"github.com/your-org/error-simulator/models"
)

type contextKey string

const ErrorTypeKey contextKey = "error_type"

// Recovery wraps an http.Handler and recovers from panics, captures stack trace,
// publishes the error event to Kafka (schema expected by ai-debugger), and returns a JSON response.
func Recovery(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				stack := string(debug.Stack())
				errMsg := panicToString(v)
				errorType := getErrorType(r)
				if errorType == "" {
					errorType = "Unknown"
				}

				event := models.ErrorEvent{
					Service:      "error-simulator",
					Repository:   cfg.GithubRepository,
					Branch:       "main",
					ErrorMessage: errMsg,
					StackTrace:   stack,
					Timestamp:    time.Now().UTC().Format(time.RFC3339),
					Environment:  "development",
				}
				kafkaErr := kafka.PublishErrorEvent(cfg, event)
				if kafkaErr != nil {
					log.Printf("[recovery] kafka publish failed: %v", kafkaErr)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error":         "panic recovered",
					"error_message": errMsg,
					"error_type":    errorType,
					"kafka_sent":    kafkaErr == nil,
					"timestamp":     time.Now().UTC().Format(time.RFC3339),
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// WithErrorType returns a handler that sets the error type in context and delegates to next.
// Recovery middleware reads this to set error_type in the HTTP response.
func WithErrorType(errorType string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := withErrorType(r.Context(), errorType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func panicToString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	if err, ok := v.(error); ok {
		return err.Error()
	}
	return fmt.Sprintf("%v", v)
}

func getErrorType(r *http.Request) string {
	v := r.Context().Value(ErrorTypeKey)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func withErrorType(ctx context.Context, errorType string) context.Context {
	return context.WithValue(ctx, ErrorTypeKey, errorType)
}
