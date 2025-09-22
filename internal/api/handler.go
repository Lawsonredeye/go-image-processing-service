package api

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png" // Import for PNG decoding side-effects
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

// CompressHandler processes an image uploaded via a multipart form and adjusts its JPEG quality.
//
// It expects a POST request with a form field named "image" containing the image file.
// The handler supports decoding of JPEG and PNG image formats.
//
// An optional query parameter `quality` (integer 1-100) can be provided.
// If the quality is not provided or is invalid, a default quality of 75 is used.
//
// The handler always returns a JPEG image.
func CompressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Could not rewind file", http.StatusInternalServerError)
		return
	}

	src, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		http.Error(w, fmt.Sprintf("Could not decode image: %v", err), http.StatusBadRequest)
		return
	}

	// Parse quality from query parameter.
	quality, err := strconv.Atoi(r.URL.Query().Get("quality"))
	if err != nil || quality < 1 || quality > 100 {
		quality = 75 // Default quality
	}

	log.Printf("Encoding with JPEG quality: %d", quality)

	opts := &jpeg.Options{Quality: quality}

	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, src, opts)
	if err != nil {
		log.Printf("Error encoding compressed image: %v", err)
		http.Error(w, "Could not encode compressed image", http.StatusInternalServerError)
		return
	}
}

// ConvertHandler processes an image and converts it to a different format.
//
// It expects a POST request with a form field named "image" containing the image file.
// A required query parameter `format` must be provided, which can be "jpeg" or "png".
//
// Upon successful processing, it returns the new image encoded in the specified format
// with the corresponding Content-Type header.
func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Could not rewind file", http.StatusInternalServerError)
		return
	}

	src, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		http.Error(w, fmt.Sprintf("Could not decode image: %v", err), http.StatusBadRequest)
		return
	}

	// Get target format from query parameter
	format := r.URL.Query().Get("format")

	switch format {
	case "jpeg", "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		err = jpeg.Encode(w, src, nil) // Use default quality
	case "png":
		w.Header().Set("Content-Type", "image/png")
		err = png.Encode(w, src)
	default:
		http.Error(w, `Invalid or missing 'format' parameter. Supported formats: jpeg, png`, http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error encoding image to %s: %v", format, err)
		http.Error(w, "Could not encode image", http.StatusInternalServerError)
	}
}

// FlipHandler processes an image and flips it horizontally or vertically.
//
// It expects a POST request with an "image" form field.
// A required `direction` query parameter must be "horizontal" or "vertical".
func FlipHandler(w http.ResponseWriter, r *http.Request) {
	// Basic boilerplate for decoding an image
	src, err := decodeImageFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	direction := r.URL.Query().Get("direction")

	var filter gift.Filter
	switch direction {
	case "horizontal":
		filter = gift.FlipHorizontal()
	case "vertical":
		filter = gift.FlipVertical()
	default:
		http.Error(w, "Invalid or missing 'direction' parameter. Supported: horizontal, vertical", http.StatusBadRequest)
		return
	}

	g := gift.New(filter)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	// Encode and send back as JPEG
	w.Header().Set("Content-Type", "image/jpeg")
	if err := jpeg.Encode(w, dst, nil); err != nil {
		http.Error(w, "Could not encode flipped image", http.StatusInternalServerError)
	}
}

// RotateHandler processes an image and rotates it by a 90-degree increment.
//
// It expects a POST request with an "image" form field.
// A required `angle` query parameter must be 90, 180, or 270.
func RotateHandler(w http.ResponseWriter, r *http.Request) {
	src, err := decodeImageFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	angle, err := strconv.Atoi(r.URL.Query().Get("angle"))
	if err != nil {
		http.Error(w, "Invalid 'angle' parameter. Must be an integer.", http.StatusBadRequest)
		return
	}

	var filter gift.Filter
	switch angle {
	case 90:
		filter = gift.Rotate90()
	case 180:
		filter = gift.Rotate180()
	case 270:
		filter = gift.Rotate270()
	default:
		http.Error(w, "Invalid 'angle' parameter. Supported: 90, 180, 270", http.StatusBadRequest)
		return
	}

	g := gift.New(filter)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	w.Header().Set("Content-Type", "image/jpeg")
	if err := jpeg.Encode(w, dst, nil); err != nil {
		http.Error(w, "Could not encode rotated image", http.StatusInternalServerError)
	}
}

// decodeImageFromRequest is a helper function to reduce boilerplate in handlers.
// It handles the request parsing, file seeking, and decoding.
func decodeImageFromRequest(r *http.Request) (image.Image, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("only POST method is allowed")
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, fmt.Errorf("failed to parse multipart form")
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("could not get uploaded file")
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("could not rewind file")
	}

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}

	return src, nil
}
