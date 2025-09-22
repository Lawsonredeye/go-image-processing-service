package api

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // Import for PNG decoding side-effects
	"log"
	"net/http"
	"strconv"

	"github.com/disintegration/gift"
)

// ResizeHandler processes an image uploaded via a multipart form and resizes it.
//
// It expects a POST request with a form field named "image" containing the image file.
// The handler supports decoding of JPEG and PNG image formats.
//
// Optional query parameters `width` and `height` (integers) can be provided to specify
// the desired dimensions. If a dimension is not provided or is invalid, it is treated as 0.
// - If both width and height are 0, a default width of 500 is used, preserving aspect ratio.
// - If one dimension is 0, it's calculated to preserve the original aspect ratio.
//
// Upon successful processing, it returns the new image encoded as a JPEG.
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

	// Parse width and height from query parameters.
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))

	// If no dimensions are provided, apply a default.
	if width == 0 && height == 0 {
		width = 500
	}

	log.Printf("Resizing to width: %d, height: %d", width, height)

	g := gift.New(
		gift.Resize(width, height, gift.LanczosResampling),
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
