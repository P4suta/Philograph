.PHONY: build test lint run clean

build:
	go build -o philograph ./cmd/philograph/

test:
	go test ./...

lint:
	go vet ./...
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"

run: build
	./philograph $(ARGS)

clean:
	rm -f philograph
	go clean -testcache
