# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run swag init
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create media directory for photos
RUN mkdir -p /root/media

# Expose port
EXPOSE 8739

# Command to run
CMD ["./main"]
