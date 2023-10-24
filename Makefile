all: build

build:
	mkdir -p bin
	go build -o bin/ ./...

test:
	go test -v ./internal/chess
	go test -v ./internal/perft/

clean:
	rm -rf bin
