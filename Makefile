.PHONY: install run-migrate-up run-migrate-down run-reset run build

install:
	go mod download

run-migrate-up:
	go run ./runner/migration/main.go up

run-migrate-down:
	go run ./runner/migration/main.go down

run-reset:
	make install && \
	go mod tidy && \
	make run-migrate-down && \
	make run-migrate-up

run:
	go run ./main.go up

build:
	go build -o bin/api ./main.go