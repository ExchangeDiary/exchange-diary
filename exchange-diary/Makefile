GO ?= GO111MODULE=on go
APP_NAME = exchange-diary
BIN_DIR = ./bin
BUILD_DIR = ./cmd
BUILD_FILE = $(addprefix $(BUILD_DIR)/, main.go)

.PHONY: build
build:
	$(GO) build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) $(BUILD_FILE)

