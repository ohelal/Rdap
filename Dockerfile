# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install git and build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o rdap-server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/rdap-server .

# Create config directory and copy config files
COPY config/ /app/config/

# Expose port
EXPOSE 8080

# Run the server
CMD ["./rdap-server"]
