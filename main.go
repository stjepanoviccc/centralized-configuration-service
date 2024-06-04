package main

import (
	"ars_projekat/config"
	"ars_projekat/handlers"
	"ars_projekat/middleware"
	"ars_projekat/repositories"
	"ars_projekat/services"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func main() {
	// Jaeger startup
	cfg := config.GetConfig()

	ctx := context.Background()
	exp, err := newExporter(cfg.JaegerAddress)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("ars_projekat")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	logger := log.New(os.Stdout, "[config-api]", log.LstdFlags)

	store, err := repositories.New(logger, tracer)
	if err != nil {
		logger.Fatal(err)
	}

	configService := services.NewConfigurationService(store, tracer)
	configHandler := handlers.NewConfigurationHandler(configService, tracer)

	configGroupService := services.NewConfigurationGroupService(store, tracer)
	configGroupHandler := handlers.NewConfigurationGroupHandler(configGroupService, tracer)

	idempotencyService := services.NewIdempotencyService(*store, tracer)
	idempotencyMiddleware := middleware.NewIdempotency(&idempotencyService, tracer)

	metricsService := services.NewMetricsService()
	metricsMiddleware := middleware.NewMetrics(metricsService)

	limiter := middleware.NewRateLimiter(time.Second, 3)

	router := mux.NewRouter()
	router.Use(otelmux.Middleware("ars_projekat"))

	router.Use(func(next http.Handler) http.Handler {
		return middleware.AdaptHandler(next, limiter)
	})
	router.Use(func(next http.Handler) http.Handler {
		return middleware.AdaptIdempotencyHandler(next, idempotencyMiddleware)
	})
	router.Use(func(next http.Handler) http.Handler {
		return middleware.AdaptPrometheusHandler(next, metricsMiddleware)
	})

	// Config routes
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	router.HandleFunc("/configs/", configHandler.Upsert).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")

	// Config group routes
	router.HandleFunc("/groups/{name}/{version}/{labels: ?.*}", configGroupHandler.Get).Methods("GET")
	router.HandleFunc("/groups/", configGroupHandler.Upsert).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}/{labels: ?.*}", configGroupHandler.Delete).Methods("DELETE")
	router.HandleFunc("/groups/{name}/{version}", configGroupHandler.AddConfig).Methods("PUT")

	// Serve the swagger.yaml file
	router.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger.yaml")
	}).Methods("GET")

	// Metrics exposer
	router.Handle("/metrics", metricsMiddleware.MetricsHandler()).Methods("GET")

	// SwaggerUI
	optionsDevelopers := swaggerMiddleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := swaggerMiddleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("Starting server..")

		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Stopped server")
}

func newExporter(address string) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ars_projekat"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
