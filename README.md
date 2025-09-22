# Go Image Processing Service (GIPS)

An image processing service built in Go, designed to handle image uploads and perform various on-the-fly transformations.

## Technologies Used

- **Language**: Go
- **Core Libraries**:
    - `net/http` for the web server.
    - `testing` and `net/http/httptest` for unit and integration tests.
    - `image`, `image/jpeg`, `image/png` for image decoding and encoding.
- **Third-Party Libraries**:
    - `github.com/disintegration/gift`: For high-quality image filtering (resize, rotate, flip).

## Features Implemented

The service exposes several endpoints for image manipulation. All endpoints expect a `POST` request with a multipart form containing an `image` field.

- **`/resize`**: Resizes an image.
    - **Query Params**: `width` (int), `height` (int)
    - **Behavior**: Preserves aspect ratio if one dimension is omitted. Uses a default width of 500px if both are omitted.
    - **Example**: `curl -X POST -F "image=@/path/to/img.png" "http://localhost:8080/resize?width=300"`

- **`/compress`**: Adjusts the quality of a JPEG image.
    - **Query Params**: `quality` (int, 1-100)
    - **Behavior**: Outputs a JPEG. Defaults to quality `75` if the parameter is missing or invalid.
    - **Example**: `curl -X POST -F "image=@/path/to/img.png" "http://localhost:8080/compress?quality=50"`

- **`/convert`**: Converts an image from one format to another.
    - **Query Params**: `format` (string, "jpeg" or "png")
    - **Behavior**: Fails if the format is missing or unsupported.
    - **Example**: `curl -X POST -F "image=@/path/to/img.jpg" "http://localhost:8080/convert?format=png"`

- **`/flip`**: Flips an image.
    - **Query Params**: `direction` (string, "horizontal" or "vertical")
    - **Behavior**: Fails if the direction is missing or invalid.
    - **Example**: `curl -X POST -F "image=@/path/to/img.png" "http://localhost:8080/flip?direction=horizontal"`

- **`/rotate`**: Rotates an image by a 90-degree increment.
    - **Query Params**: `angle` (int, 90, 180, or 270)
    - **Behavior**: Fails if the angle is missing or unsupported.
    - **Example**: `curl -X POST -F "image=@/path/to/img.png" "http://localhost:8080/rotate?angle=90"`

- **`/crop`**: Crops an image to a specified rectangle.
    - **Query Params**: `x` (int), `y` (int), `width` (int), `height` (int)
    - **Behavior**: Fails if any parameter is missing or invalid.
    - **Example**: `curl -X POST -F "image=@/path/to/img.png" "http://localhost:8080/crop?x=10&y=10&width=100&height=100"`

## Setup and Run Instructions

### Prerequisites
- Go (version 1.18 or later)

### Running the Service
1.  Clone the repository.
2.  Start the server:
    ```sh
    go run ./cmd/api
    ```
    The service will be available at `http://localhost:8080`.

### Running Tests
To run the complete test suite:
```sh
go test ./...
```

## Notes on AI Usage

This project was developed with the assistance of a large language model (Gemini). The AI was used in the following contexts:

- **Initial Scaffolding**: Generated the initial project structure (`cmd`/`internal`) and boilerplate code.
- **Test-Driven Development (TDD)**: For each feature, the AI was prompted to write a comprehensive, table-driven test suite first. These tests initially failed as expected.
- **Feature Implementation**: The AI then generated the handler logic required to make the tests pass for each feature (`resize`, `compress`, `convert`, `flip`, `rotate`).
- **Debugging**: Collaboratively debugged issues, including multipart form decoding errors (`file.Seek`), missing image format decoders (`image: unknown format`), and Go syntax errors (`imported and not used`).
- **Refactoring and Documentation**: Refactored the code to improve structure (e.g., adding the `decodeImageFromRequest` helper) and generated clear, well-formatted function documentation and this README file.
