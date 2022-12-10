all: build

build:
	mkdir -p bin
	go build -o bin/ ./...

test:
	go test -v ./internal/chess

clean:
	rm -rf bin
