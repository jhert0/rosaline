all: build

build:
	go build -o bin/rosaline ./cmd/rosaline/
	go build -o bin/playground ./cmd/playground/

test:
	go test -v ./internal/chess

clean:
	rm -rf bin
