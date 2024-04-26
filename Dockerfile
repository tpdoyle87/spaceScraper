# Use the official Golang image from the Docker Hub.
FROM golang:1.21.3

# Set the working directory inside the container.
WORKDIR /app

# Copy the local code to the container's workspace.
COPY . .

# Install the dependencies specified in go.mod and go.sum.
RUN go mod download

# Install PostgreSQL and other system dependencies (assuming Debian-based image).
RUN apt-get update && apt-get install -y postgresql-client

# Build the application.
RUN go build -o main .

# Expose port 8080 for web access.
EXPOSE 8080

# Command to run the executable.
CMD ["./main"]
