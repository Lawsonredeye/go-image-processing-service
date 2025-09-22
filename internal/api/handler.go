package api

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Import for PNG decoding side-effects
	"log"
	"net/http"

	"github.com/disintegration/gift"
)

// ResizeHandler processes an image uploaded via a multipart form.
//
// It expects a POST request with a form field named "image" containing the image file.
// The handler supports decoding of JPEG and PNG image formats.
//
// Upon successful processing, it performs a resize operation (currently to a fixed
// width of 500px while preserving aspect ratio) and returns the new image
// encoded as a JPEG with a "Content-Type" of "image/jpeg".
//
// If any step fails (e.g., missing file, decoding error), it returns an appropriate
// HTTP error status code and a descriptive error message.
func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Received file: %s, size: %d bytes, content-type: %s", header.Filename, header.Size, header.Header.Get("Content-Type"))

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Could not rewind file", http.StatusInternalServerError)
		return
	}

	src, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		http.Error(w, fmt.Sprintf("Could not decode image: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("Successfully decoded image, format: %s", format)

	g := gift.New(
		gift.Resize(500, 0, gift.LanczosResampling),
	)

	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, dst, nil)
	if err != nil {
		log.Printf("Error encoding resized image: %v", err)
		http.Error(w, "Could not encode resized image", http.StatusInternalServerError)
		return
	}
	log.Println("Successfully resized and sent image.")
}
