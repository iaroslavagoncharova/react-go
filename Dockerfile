# Use the official Golang image as a base
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and install dependencies first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which the app will run
EXPOSE 8080

# Command to run the executable
CMD ["./main"]