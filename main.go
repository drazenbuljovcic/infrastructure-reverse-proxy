package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joho/godotenv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(os.Getenv("OTEL_SERVICE_NAME")),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

func logRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s %s", r.Method, r.URL.Host, r.URL.Path, )
		next.ServeHTTP(w, r)
	})
}

func reverseProxyHandlerMiddleware(targetURL *url.URL) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		r.Host = targetURL.Host
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme

		// Extract traceparent from the incoming request
		ctx := r.Context()
		traceparent := r.Header.Get("traceparent")
		log.Printf("traceparent: %s", traceparent)
		if traceparent != "" {
			ctx = propagation.TraceContext{}.Extract(ctx, propagation.HeaderCarrier(r.Header))
		}

		// start a new span for the actual proxy request and defer its end
		_, span := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME")).Start(ctx, r.URL.Path)
		defer span.End()
		
		childTraceparent := fmt.Sprintf("00-%s-%s-01", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		// for key, values := range r.Header {
		// 	for _, value := range values {
		// 		r.Header.Add(key, value)
		// 	}
		// }

		r.Header.Del("traceparent")
		r.Header.Add("traceparent", childTraceparent)

		// Perform the reverse proxy
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	// Replace "http://backend-server.com" with the actual backend server URL
	// backendURL, err := url.Parse("http://172.25.0.3:3000")
	reverseProxyBackendHost := os.Getenv("SERVICE_HOSTNAME")

	// write traces to a file?
	l := log.New(os.Stdout, "", 0)
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	// set up the exporter
	// exp, err := newExporter(f)
	// if err != nil {
	// 	l.Fatal(err)
	// }

	exporter, err := zipkin.New(
		os.Getenv("OTEL_API_HOST") + "/api/v2/spans",
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(newResource()),
	)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)
	
	backendURL, err := url.Parse(reverseProxyBackendHost)
	if err != nil {
		log.Fatal("Error parsing backend URL:", err)
	}

	// Create the reverse proxy handler
	proxyHandler := reverseProxyHandlerMiddleware(backendURL)

	// Use the logging middleware to wrap the reverse proxy handler
	handlerWithLogging := logRequestsMiddleware(proxyHandler)

	// Define a handler function for the "/" route
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// Update the request's Host header to match the backend server's hostname
	// 	r.Host = backendURL.Host

	// 	// Serve the request using the reverse proxy
	// 	proxyHandler.ServeHTTP(w, r)
	// })

	// Start the HTTP server on port 8080
	// log.Printf("Starting reverse proxy server on port %s...\n", port)
	log.Printf("Started!! Go!")

	log.Printf("Application started on port %s!", os.Getenv("PORT"))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), handlerWithLogging); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
