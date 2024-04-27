# Use the official Golang image to create a build artifact.
FROM golang:1.21.3 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files from the root directory if necessary
# Update the path according to where go.mod and go.sum are located
COPY ../../go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the entire project and build it
# This assumes the entire project is needed to build the binary
COPY ../.. .

# Build the Go app
# Change to the directory where the main Go file is located before building
WORKDIR /app/cmd/scraper
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o spaceScraper .

# Start a new stage from scratch for a cleaner, smaller image
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/cmd/scraper/spaceScraper .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./spaceScraper"]
