SERVER_BIN=go-dns-server
SERVER_SRC=cmd/server/main.go
SERVER_TEST_SRC=cmd/server/main_test.go
SERVER_PORT=7777

.PHONY: run clean all test

build:
	go build -o $(SERVER_BIN) $(SERVER_SRC)

run: build
	@echo "Starting server..."
	@./$(SERVER_BIN) --port $(SERVER_PORT) &

test: build run
	@echo "Running tests..."
	go test -v $(SERVER_TEST_SRC)
	@make kill clean

clean:
	@echo "Cleaning up..."
	@rm -f $(SERVER_BIN)

kill: 
	@echo "Stopping server..."
	@kill -9 $$(lsof -t -i udp:$(SERVER_PORT))

all: build run kill clean
