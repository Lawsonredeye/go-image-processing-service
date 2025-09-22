
package main

import (
	"go-image-processing-service/internal/server"
)

// main is the entry point for the image processing service.
// It creates and starts a new server instance.
func main() {
	srv := server.New("8080")
	srv.Start()
}
