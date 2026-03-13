package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/your-org/error-simulator/config"
	"github.com/your-org/error-simulator/logger"

	"go.elastic.co/apm/v2"
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

				// Report to Elastic APM for Kibana APM Services
				apm.CaptureError(r.Context(), errors.New(errMsg)).Send()

				// Structured log so we know exactly what fucked up when raising fix PRs.
				// Include repository/branch so Logstash→Kafka events have them for ai-debugger PR creation.
				logger.Log.Error().
					Str("error_type", errorType).
					Str("path", r.URL.Path).
					Str("method", r.Method).
					Str("error_message", errMsg).
					Str("stack_trace", stack).
					Str("what_failed", whatFailedSummary(errorType, r.URL.Path)).
					Str("service", "error-simulator").
					Str("repository", cfg.GithubRepository).
					Str("branch", "main").
					Msg("panic recovered")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error":         "panic recovered",
					"error_message": errMsg,
					"error_type":    errorType,
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

// whatFailedSummary returns a one-liner describing what to fix (for PR titles/descriptions).
func whatFailedSummary(errorType, path string) string {
	m := map[string]string{
		"NilPointer":    "order.Patient nil dereference in OrderService.ProcessOrder",
		"DBError":       "nil *sql.DB in UserRepository.GetUserByID",
		"Panic":         "PaymentService.ProcessPayment panic when amount exceeds limit",
		"IndexOOB":      "ReportGenerator.GetTopProducts index out of range (slice len < 6)",
		"TypeAssertion": "ConfigLoader.GetDatabaseConfig type assertion on config[\"database\"]",
		"DivisionZero":  "MetricsService.CalculateConversionRate divide by zero (totalVisits=0)",
		"Deadlock":      "CacheManager mutex ordering (UpdateCache vs InvalidateCache)",
		"StackOverflow": "TreeNode.CalculateDepth missing nil base case / circular ref",
	}
	if s, ok := m[errorType]; ok {
		return s
	}
	return path + " — unknown error type " + errorType
}
