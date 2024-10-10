FROM golang:1.21 AS builder

# enable Go modules support
ENV GO111MODULE=on

WORKDIR $GOPATH/src/privy-middleware

# manage dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download


# Install migrate tool
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate.linux-amd64 /usr/local/bin/migrate

# Copy src code from the host and compile it
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /privy main.go

CMD ["/bin/myapp"]

FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=builder privy .
COPY migration migration



# Expose the port your Go application listens on
EXPOSE 9001

# Command to run the application
CMD ["./privy"]
