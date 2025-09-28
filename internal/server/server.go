package server

import (
	"fmt"
	"log"
	"net/http"

	"go-image-processing-service/internal/api"
)

// Server holds the dependencies and configuration for our HTTP server.
type Server struct {
	port string
}

// New creates and returns a new Server instance, configured to listen on the given port.
func New(port string) *Server {
	return &Server{
		port: port,
	}
}

// Start initializes all server routes and begins listening for incoming HTTP requests.
// It will block until the server is stopped or a fatal error occurs.
func (s *Server) Start() {
	// Create a new mux (router)
	rootMux := http.NewServeMux()
	mux := http.NewServeMux()

	mux.HandleFunc("/health", api.HealthCheckHandler)
	mux.HandleFunc("/resize", api.ResizeHandler)
	mux.HandleFunc("/compress", api.CompressHandler)
	mux.HandleFunc("/convert", api.ConvertHandler)
	mux.HandleFunc("/flip", api.FlipHandler)
	mux.HandleFunc("/rotate", api.RotateHandler)
	mux.HandleFunc("/crop", api.CropHandler)

	rootMux.Handle("/api/", http.StripPrefix("/api", mux))

	// Wrap the mux with a CORS middleware
	h := corsMiddleware(rootMux)

	fmt.Printf("Starting server on http://localhost:%s\n", s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, h))
}

// corsMiddleware is a simple middleware to handle CORS.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For development, we can allow any origin.
		// For production, you would want to restrict this to your frontend's domain.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle pre-flight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
