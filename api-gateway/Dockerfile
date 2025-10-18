# ---------- Build Stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (some Go modules may need it)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build with optimizations for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gateway main.go


# ---------- Runtime Stage ----------
FROM alpine:latest

# Install certificates (needed if your app calls HTTPS endpoints)
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/gateway .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./gateway"]
