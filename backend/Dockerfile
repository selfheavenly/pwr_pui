# Build stage: compile the Go binary
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary and output as "main"
RUN go build -o main .

# Run stage: use a minimal container image
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8000

# Run the binary
CMD ["./main"]
