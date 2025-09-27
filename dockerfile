# Stage 1: Builder
# This stage uses the Go image to build the application.
FROM golang:1.25 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
# This leverages Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Build the Go application as a static binary
# CGO_ENABLED=0 is important for creating a static binary that can run in a scratch image.
# The output is a single binary file named 'server' in the /app directory.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/server ./cmd/api

# Stage 2: Final
# This stage creates the final, small production image.
FROM scratch

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/server /server

# Expose the port the application runs on
EXPOSE 8080

# Set the command to run the executable
CMD ["/server"]
