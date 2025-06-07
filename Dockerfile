# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates imagemagick

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create media directory for photos
RUN mkdir -p /root/media

# Command to run
CMD ["./main"]
