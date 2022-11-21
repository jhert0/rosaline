all: build

build:
	go build -o bin/rosaline ./cmd/rosaline/

test:
	go test -v ./internal/chess

clean:
	rm -rf bin
