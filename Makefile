SERVER_BIN=go-dns-server
SERVER_SRC=cmd/server/main.go
SERVER_TEST_SRC=cmd/server/main_test.go

.PHONY: run clean all test

build:
	go build -o $(SERVER_BIN) $(SERVER_SRC)

run: build
	@echo "Starting server..."
	./$(SERVER_BIN) &

test: build
	@echo "Starting server..."
	./$(SERVER_BIN) &

	@go test -v $(SERVER_TEST_SRC)

	@make kill
	@make clean

clean:
	@echo "Cleaning up..."
	@rm -f $(SERVER_BIN)

kill: 
	@echo "Stopping server..."
	@kill -9 $$(lsof -t -i udp:2053)

all: build run kill clean
