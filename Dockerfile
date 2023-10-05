# Use the official Go image as the base image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /project-privy

# Copy the Go project source code into the container
COPY . .

# Build the Go application
RUN go build


# Stage 2: Create the PostgreSQL database
FROM postgres:14 AS database

# Environment variables for PostgreSQL
ENV POSTGRES_DB=privy
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=p@$$w0rdprivy

# Copy the SQL script to create the database
COPY init.sql /docker-entrypoint-initdb.d/

# Stage 3: Final image
FROM ubuntu:20.04

# Install any additional dependencies if needed

# Copy the Go binary from the builder stage
COPY --from=builder /app/your-app-binary /app/

# Copy any other files your application needs

# Copy the PostgreSQL configuration and data from the database stage
COPY --from=database /var/lib/postgresql /var/lib/postgresql
COPY --from=database /etc/postgresql /etc/postgresql

# Expose the port your Go application listens on
EXPOSE 8080

# Command to run the application
CMD ["./project-privy"]