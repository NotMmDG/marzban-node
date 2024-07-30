# Build Stage
FROM golang:1.18-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Run Stage
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file from the build stage
COPY --from=build /app/main .

# Expose the application on port 62050 (replace with the appropriate port)
EXPOSE 62050

# Command to run the executable
CMD ["./main"]
