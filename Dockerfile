FROM alpine:latest

#RUN apt-get update && apt-get install -y golang-go ca-certificates
# Set the working directory inside the container
# WORKDIR /project-privy

# Copy the Go project source code into the container
COPY privy .

# Expose the port your Go application listens on
EXPOSE 9001

# Command to run the application
CMD ["./privy"]
