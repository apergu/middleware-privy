FROM golang:1.21 AS builder

# enable Go modules support
ENV GO111MODULE=on

WORKDIR $GOPATH/src/privy-middleware

# manage dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy src code from the host and compile it
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /privy main.go

CMD ["/bin/myapp"]


FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=builder /privy /bin
#RUN apt-get update && apt-get install -y golang-go ca-certificates
# Set the working directory inside the container
# WORKDIR /project-privy

# Copy the Go project source code into the container
# COPY privy .

# Expose the port your Go application listens on
EXPOSE 9001

# Command to run the application
CMD ["./bin/privy"]
