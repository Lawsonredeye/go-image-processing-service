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
