# Use the official Go image as the base image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project source code into the container
COPY . .

# Build the Go application
RUN go build

# Expose the port your Go application listens on
EXPOSE 8080

# Command to run the application
CMD ["./project-privy"]