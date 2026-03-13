package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-org/error-simulator/cachesvc"
	"github.com/your-org/error-simulator/config"
	"github.com/your-org/error-simulator/configsvc"
	"github.com/your-org/error-simulator/handlers"
	"github.com/your-org/error-simulator/kafka"
	"github.com/your-org/error-simulator/logger"
	"github.com/your-org/error-simulator/middleware"
	"github.com/your-org/error-simulator/pipeline"
	"github.com/your-org/error-simulator/userfetcher"
	"github.com/your-org/error-simulator/usersvc"
)

func main() {
	cfg := config.Load()
	kafka.InitProducer(cfg)

	orderSvc := handlers.NewOrderService()
	userRepo := handlers.NewUserRepository()
	paymentSvc := handlers.NewPaymentService()
	reportGen := handlers.NewReportGenerator()
	metricsSvc := handlers.NewMetricsService("monthly")
	cacheMgr := handlers.NewCacheManager()
	treeOps := &handlers.TreeOps{}
	orderPipeline := &pipeline.Pipeline{}
	configSvc := &configsvc.Service{}
	cacheSvcRepo := cachesvc.NewRepo()
	cacheSvc := cachesvc.NewCacheService(cacheSvcRepo)
	userFetcherImpl := userfetcher.NewImpl()
	userSvc := &usersvc.Service{Fetcher: userFetcherImpl}

	// WithErrorType must wrap Recovery so that when we recover, r.Context() has the error type.
	wrap := func(errorType string, h http.Handler) http.Handler {
		return middleware.WithErrorType(errorType, middleware.Recovery(cfg, h))
	}

	mux := http.NewServeMux()
	mux.Handle("/error/nil-pointer", wrap("NilPointer", handlers.NilPointer(orderSvc)))
	mux.Handle("/error/db", wrap("DBError", handlers.DBError(userRepo)))
	mux.Handle("/error/panic", wrap("Panic", handlers.PanicRecovery(paymentSvc)))
	mux.Handle("/error/index-oob", wrap("IndexOOB", handlers.IndexOOB(reportGen)))
	mux.Handle("/error/type-assertion", wrap("TypeAssertion", handlers.TypeAssertion(nil))) // loader created inside handler
	mux.Handle("/error/division-zero", wrap("DivisionZero", handlers.DivisionZero(metricsSvc)))
	mux.Handle("/error/deadlock", wrap("Deadlock", handlers.Deadlock(cacheMgr)))
	mux.Handle("/error/stack-overflow", wrap("StackOverflow", handlers.StackOverflow(treeOps)))
	mux.Handle("/error/multi-file/order", wrap("MultiFileOrder", handlers.MultiFileOrder(orderPipeline)))
	mux.Handle("/error/multi-file/config", wrap("MultiFileConfig", handlers.MultiFileConfig(configSvc)))
	mux.Handle("/error/multi-file/cache", wrap("MultiFileCache", handlers.MultiFileCache(cacheSvc)))
	mux.Handle("/error/multi-file/interface", wrap("MultiFileInterface", handlers.MultiFileInterface(userSvc)))
	mux.Handle("/error/multi-file/callback", wrap("MultiFileCallback", handlers.MultiFileCallback()))
	mux.Handle("/error/latency", handlers.Latency(cfg))

	printBanner(cfg)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("server listen failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info().Msg("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error().Err(err).Msg("shutdown error")
	}
	logger.Log.Info().Msg("server stopped")
}

func printBanner(cfg *config.Config) {
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	kafkaAddr := cfg.KafkaBootstrapServers
	if kafkaAddr == "" {
		kafkaAddr = "(not set)"
	}
	topic := cfg.KafkaTopic
	if topic == "" {
		topic = "service.errors"
	}
	fmt.Println("=== Error Simulator Server ===")
	fmt.Printf("Port: %s\n\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET /error/nil-pointer      → nil pointer dereference (order.Patient.Name)")
	fmt.Println("  GET /error/db               → nil db connection panic (UserRepository.GetUserByID)")
	fmt.Println("  GET /error/panic            → explicit panic (PaymentService.ProcessPayment)")
	fmt.Println("  GET /error/index-oob        → index out of range (ReportGenerator.GetTopProducts)")
	fmt.Println("  GET /error/type-assertion   → type assertion failure (ConfigLoader.GetDatabaseConfig)")
	fmt.Println("  GET /error/division-zero    → integer divide by zero (MetricsService.CalculateConversionRate)")
	fmt.Println("  GET /error/deadlock         → goroutine deadlock (CacheManager)")
	fmt.Println("  GET /error/stack-overflow   → stack overflow (TreeNode.CalculateDepth)")
	fmt.Println("  GET /error/multi-file/order → 3-file: handler→pipeline→formatter (nil ShippingAddress)")
	fmt.Println("  GET /error/multi-file/config   → 3-file: handler→configsvc→env (type assertion)")
	fmt.Println("  GET /error/multi-file/cache    → 3-file: handler→cachesvc→repo (nil db)")
	fmt.Println("  GET /error/multi-file/interface → genre: interface (handler→usersvc→userfetcher impl panic)")
	fmt.Println("  GET /error/multi-file/callback  → genre: callback (handler callback panics via processor)")
	fmt.Printf("  GET /error/latency              → configurable delay (default %d ms, override ?ms=)\n", cfg.LatencyMs)
	fmt.Printf("\nKafka: %s → %s\n", kafkaAddr, topic)
	fmt.Println("==============================")
}
