all: build

build:
	mkdir -p bin
	go build -o bin/ ./...

test:
	go test -v ./internal/chess

perft-test:
	go test -v ./internal/perft/

test-all: test perft-test

clean:
	rm -rf bin

coverage:
	go test -cover ./internal/chess
