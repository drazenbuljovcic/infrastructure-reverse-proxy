package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func logRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s %s", r.Method, r.URL.Host, r.URL.Path, )
		next.ServeHTTP(w, r)
	})
}

func reverseProxyHandlerMiddleware(targetURL *url.URL) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		// Modify request headers if needed
		r.Host = targetURL.Host
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme

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

	log.Print("HOSTNAME: ", reverseProxyBackendHost)
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
