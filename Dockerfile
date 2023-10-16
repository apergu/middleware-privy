# Use the official Go image as the base image
FROM golang:1.18-alpine

#RUN apt-get update && apt-get install -y golang-go ca-certificates
# Set the working directory inside the container
WORKDIR /project-privy

# Copy the Go project source code into the container
COPY . .

# Configure Go
RUN go install

# Build the Go application
RUN go build -v -o project-privy .

# Expose the port your Go application listens on
EXPOSE 9001

# Command to run the application
CMD ["./project-privy"]
