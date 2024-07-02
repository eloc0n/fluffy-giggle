# syntax=docker/dockerfile:1

# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Install Git
RUN apk update && \
    apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Stage 2: Run the Go application with Air
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Build the Go app

# Copy the Go modules manifests and download dependencies
COPY --from=builder /go/bin/air /usr/bin/air

# Copy the application source code
COPY . .

# Copy the Air configuration file
COPY .air.toml .

RUN go build -o ./tmp/main ./api/main.go

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application with Air
CMD ["air"]
