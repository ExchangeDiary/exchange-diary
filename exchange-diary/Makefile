GO ?= GO111MODULE=on go
APP_NAME = exchange-diary
BIN_DIR = ./bin
BUILD_DIR = ./cmd
BUILD_FILE = $(addprefix $(BUILD_DIR)/, main.go)

.PHONY: build
build: $(addprefix $(BIN_DIR)/, $(APP_NAME)) # ./bin/exchange-diary

$(BIN_DIR)/%: $(shell find . | grep $(APP_NAME))
	$(GO) build -ldflags="-s -w" -o $@ $(BUILD_FILE)

