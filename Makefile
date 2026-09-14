BINARY   := parameters-kerberos
CMD_PKG  := ./cmd/server
BUILD_DIR := dist

.PHONY: all build test lint clean docker run

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD_PKG)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)

docker:
	docker build -t $(BINARY):latest .

run: build
	$(BUILD_DIR)/$(BINARY)

tidy:
	go mod tidy

vet:
	go vet ./...
