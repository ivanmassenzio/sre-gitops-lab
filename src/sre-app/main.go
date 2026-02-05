package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer    trace.Tracer
	errorRate int
	latencyMs int
)

// Metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func initTracer() func(context.Context) error {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("observability-tempo.monitoring.svc.cluster.local:4317"), // Direct to Tempo/Collector
	)
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("sre-observability-app"),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "lab"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer = tp.Tracer("sre-observability-app")

	return tp.Shutdown
}

func main() {
	shutdown := initTracer()
	defer shutdown(context.Background())

	// Env configs
	errorRate, _ = strconv.Atoi(os.Getenv("ERROR_RATE")) // 0-100
	latencyMs, _ = strconv.Atoi(os.Getenv("LATENCY_MS")) // milliseconds

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(handleRoot), "root"))
	mux.Handle("/checkout", otelhttp.NewHandler(http.HandlerFunc(handleCheckout), "checkout"))

	log.Println("Starting SRE App on :8080")
	log.Printf("Config: ERROR_RATE=%d%%, LATENCY_MS=%dms\n", errorRate, latencyMs)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx, span := tracer.Start(r.Context(), "handleRoot")
	defer span.End()

	simulateWork(ctx)

	status := http.StatusOK
	if shouldError() {
		status = http.StatusInternalServerError
		span.SetAttributes(attribute.Bool("error", true))
		span.RecordError(fmt.Errorf("artificial chaos error"))
		http.Error(w, "Chaos Monkey struck!", status)
		log.Printf("Error injected 500")
	} else {
		fmt.Fprintf(w, "Hello from SRE App! TraceID: %s\n", span.SpanContext().TraceID().String())
	}

	duration := time.Since(start).Seconds()
	httpRequestsTotal.WithLabelValues("/", strconv.Itoa(status)).Inc()
	httpRequestDuration.WithLabelValues("/").Observe(duration)
}

func handleCheckout(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx, span := tracer.Start(r.Context(), "handleCheckout")
	defer span.End()

	// Simulate a database call
	dbCtx, dbSpan := tracer.Start(ctx, "database_query")
	time.Sleep(time.Duration(20+rand.Intn(50)) * time.Millisecond)
	dbSpan.SetAttributes(attribute.String("db.system", "postgres"), attribute.String("db.statement", "SELECT * FROM cart"))
	dbSpan.End()

	simulateWork(dbCtx)

	status := http.StatusOK
	if shouldError() {
		status = http.StatusInternalServerError
		http.Error(w, "Checkout failed", status)
	} else {
		fmt.Fprintf(w, "Checkout successful")
	}

	duration := time.Since(start).Seconds()
	httpRequestsTotal.WithLabelValues("/checkout", strconv.Itoa(status)).Inc()
	httpRequestDuration.WithLabelValues("/checkout").Observe(duration)
}

func simulateWork(ctx context.Context) {
	_, span := tracer.Start(ctx, "simulateWork")
	defer span.End()

	if latencyMs > 0 {
		time.Sleep(time.Duration(latencyMs) * time.Millisecond)
		span.SetAttributes(attribute.Int("simulated_latency_ms", latencyMs))
	}
}

func shouldError() bool {
	if errorRate <= 0 {
		return false
	}
	return rand.Intn(100) < errorRate
}
