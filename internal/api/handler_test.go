package api_test

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-image-processing-service/internal/api"
)

// createDummyImage generates a 10x10 PNG image in memory for testing.
func createDummyImage() (*bytes.Buffer, error) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	return buf, err
}

// createImageUploadRequest creates a new multipart/form-data HTTP request with a dummy image.
func createImageUploadRequest(url string, body io.Reader, contentType string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	return req
}

func TestResizeHandler(t *testing.T) {
	// --- Test Cases Definition ---
	testCases := []struct {
		name               string
		requestSetup       func() *http.Request
		expectedStatusCode int
		checkResponse      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success - Custom Width and Height",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/resize?width=5&height=5", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				img, _, err := image.Decode(recorder.Body)
				if err != nil {
					t.Fatalf("Failed to decode response image: %v", err)
				}
				if img.Bounds().Dx() != 5 || img.Bounds().Dy() != 5 {
					t.Errorf("Expected image dimensions 5x5, got %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
				}
			},
		},
		{
			name: "Success - Default Resize",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/resize", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				img, _, err := image.Decode(recorder.Body)
				if err != nil {
					t.Fatalf("Failed to decode response image: %v", err)
				}
				if img.Bounds().Dx() != 500 {
					t.Errorf("Expected image width 500, got %d", img.Bounds().Dx())
				}
			},
		},
		{
			name: "Failure - GET Method Not Allowed",
			requestSetup: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/resize", nil)
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "Failure - No File Uploaded",
			requestSetup: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/resize", nil)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failure - Non-Image File",
			requestSetup: func() *http.Request {
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.txt")
				part.Write([]byte("this is not an image"))
				writer.Close()
				return createImageUploadRequest("/resize", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	// --- Test Runner ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := tc.requestSetup()

			api.ResizeHandler(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
				body, _ := io.ReadAll(recorder.Body)
				t.Logf("Response body: %s", string(body))
			}

			if tc.checkResponse != nil {
				tc.checkResponse(t, recorder)
			}
		})
	}
}

func TestCompressHandler(t *testing.T) {
	// --- Test Cases Definition ---
	testCases := []struct {
		name               string
		requestSetup       func() *http.Request
		expectedStatusCode int
	}{
		{
			name: "Success - Custom Quality",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/compress?quality=50", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Success - Default Quality",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/compress", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Success - Invalid Quality Fallback to Default",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/compress?quality=abc", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	// --- Test Runner ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := tc.requestSetup()

			// This will fail to compile until we create the handler.
			api.CompressHandler(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestConvertHandler(t *testing.T) {
	// --- Test Cases Definition ---
	testCases := []struct {
		name               string
		requestSetup       func() *http.Request
		expectedStatusCode int
		expectedMimeType   string
	}{
		{
			name: "Success - PNG to JPEG",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/convert?format=jpeg", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusOK,
			expectedMimeType:   "image/jpeg",
		},
		{
			name: "Success - Upload JPEG, Convert to PNG",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage() // We can still upload a PNG, decode is agnostic
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/convert?format=png", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusOK,
			expectedMimeType:   "image/png",
		},
		{
			name: "Failure - Missing Format",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/convert", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid Format",
			requestSetup: func() *http.Request {
				imgBuf, _ := createDummyImage()
				body := new(bytes.Buffer)
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("image", "test.png")
				part.Write(imgBuf.Bytes())
				writer.Close()
				return createImageUploadRequest("/convert?format=gif", body, writer.FormDataContentType())
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	// --- Test Runner ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := tc.requestSetup()

			// This will fail to compile until we create the handler.
			api.ConvertHandler(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}

			if recorder.Code == http.StatusOK {
				contentType := recorder.Header().Get("Content-Type")
				if contentType != tc.expectedMimeType {
					t.Errorf("Expected Content-Type %s, got %s", tc.expectedMimeType, contentType)
				}
			}
		})
	}
}

func TestFlipHandler(t *testing.T) {
	testCases := []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"Success - Horizontal", "/flip?direction=horizontal", http.StatusOK},
		{"Success - Vertical", "/flip?direction=vertical", http.StatusOK},
		{"Failure - Missing Direction", "/flip", http.StatusBadRequest},
		{"Failure - Invalid Direction", "/flip?direction=diagonal", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imgBuf, _ := createDummyImage()
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("image", "test.png")
			part.Write(imgBuf.Bytes())
			writer.Close()
			req := createImageUploadRequest(tc.url, body, writer.FormDataContentType())

			recorder := httptest.NewRecorder()
			api.FlipHandler(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestRotateHandler(t *testing.T) {
	testCases := []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"Success - 90 degrees", "/rotate?angle=90", http.StatusOK},
		{"Success - 180 degrees", "/rotate?angle=180", http.StatusOK},
		{"Failure - Missing Angle", "/rotate", http.StatusBadRequest},
		{"Failure - Invalid Angle", "/rotate?angle=45", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imgBuf, _ := createDummyImage()
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("image", "test.png")
			part.Write(imgBuf.Bytes())
			writer.Close()
			req := createImageUploadRequest(tc.url, body, writer.FormDataContentType())

			recorder := httptest.NewRecorder()
			api.RotateHandler(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestCropHandler(t *testing.T) {
	testCases := []struct {
		name               string
		url                string
		expectedStatusCode int
		checkResponse      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success - Valid Crop",
			url:  "/crop?x=2&y=2&width=5&height=5",
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				img, _, err := image.Decode(recorder.Body)
				if err != nil {
					t.Fatalf("Failed to decode response image: %v", err)
				}
				if img.Bounds().Dx() != 5 || img.Bounds().Dy() != 5 {
					t.Errorf("Expected image dimensions 5x5, got %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
				}
			},
		},
		{
			name: "Failure - Missing All Params",
			url:  "/crop",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failure - Missing One Param",
			url:  "/crop?x=10&y=10&width=50", // Missing height
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid Param Type",
			url:  "/crop?x=10&y=ten&width=50&height=50",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imgBuf, _ := createDummyImage()
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("image", "test.png")
			part.Write(imgBuf.Bytes())
			writer.Close()
			req := createImageUploadRequest(tc.url, body, writer.FormDataContentType())

			recorder := httptest.NewRecorder()
			api.CropHandler(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}

			if tc.checkResponse != nil {
				tc.checkResponse(t, recorder)
			}
		})
	}
}
