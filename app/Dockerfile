# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy application source
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage - use minimal base image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# Copy the binary from builder
COPY --from=builder /app/main .

# Use nonroot user
USER nonroot:nonroot

# Cloud Run will set PORT environment variable
EXPOSE 8080

# Run the application
ENTRYPOINT ["/main"]
