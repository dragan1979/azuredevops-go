# --- Stage 1: Build Stage ---
FROM golang:1.25.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the module definition files first
# This allows Docker to cache the dependency download step
COPY go.mod go.sum ./


# Download all dependencies defined in go.mod

RUN go mod download

# Copy the rest of the source code
COPY main.go db.go handlers.go ./
COPY templates ./templates

RUN go mod tidy

# Build the final application binary
# CGO_ENABLED=0 creates a statically linked binary, making it easy to run
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-manager .


# --- STAGE 2: Final Stage (Uses a minimal environment for small image size) ---
FROM alpine:latest

# Install necessary runtime dependencies (like time zone data for Go's time package)
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory to /app
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/user-manager .
COPY --from=builder /app/templates /app/templates

RUN ls -la /app/ && ls -la /app/templates/

# Command to run the application when the container starts
CMD ["./user-manager"]