all: build

build:
	go build -o bin/rosaline ./cmd/rosaline/

test:
	go test ./internal/chess
