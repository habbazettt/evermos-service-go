# Use official Golang image as base
FROM golang:1.20 AS builder

# Set working directory
WORKDIR /app

# Copy all files
COPY . .

# Download dependencies
RUN go mod tidy

# Build the application
RUN go build -o server main.go

# Create a minimal runtime image
FROM gcr.io/distroless/base-debian11

# Set working directory
WORKDIR /root/

# Copy compiled binary
COPY --from=builder /app/server .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./server"]
