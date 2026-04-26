# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /oauth-server ./cmd/api/main.go

# Runtime stage
FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /oauth-server .

# Copy web assets for the UI
COPY --from=builder /app/web ./web

# Expose the application port
EXPOSE 8080

# Run the server
CMD ["./oauth-server"]
