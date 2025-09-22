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
	http.HandleFunc("/resize", api.ResizeHandler)
	http.HandleFunc("/compress", api.CompressHandler)
	http.HandleFunc("/convert", api.ConvertHandler)

	fmt.Printf("Starting server on http://localhost:%s\n", s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, nil))
}
