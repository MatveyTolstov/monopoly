.PHONY: build run test fmt vet clean

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -rf bin/
