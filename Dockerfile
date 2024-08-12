# Step 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Cache dependencies by copying go.mod and go.sum
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum aren't changed
RUN go mod download

# Copy the source code into the container
COPY cmd cmd
COPY internal internal
COPY *.go ./

# Build the Go app
RUN go build -o myapp

# Step 2: Create a smaller image and copy the binary
FROM alpine:3.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/myapp .

# Command to run the executable
ENTRYPOINT ["./myapp", "serve"]
