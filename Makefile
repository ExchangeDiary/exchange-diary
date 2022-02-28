.PHONY: run build docker up down clean

# export CGO_ENABLED=0
# export GOOS=linux
# export GOARCH=amd64

GO ?= GO111MODULE=on go
APP_NAME = exchange-diary
BIN_DIR = ./bin
BUILD_DIR = ./application/cmd
BUILD_FILE = $(addprefix $(BUILD_DIR)/, main.go)

# local run
run:
	$(GO) run $(BUILD_FILE)

# build binary
build:
	$(GO) build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) $(BUILD_FILE)

docker:
	make dbuild && make drun

# docker build
dbuild:
	docker build \
		-t $(APP_NAME):latest \
		-f Dockerfile --no-cache .

# docker local run
drun:
	docker run --rm -p 8080:8080 --name exchange-diary exchange-diary

# docker compose up
up:
	docker compose up -d --build --remove-orphans
	docker compose logs -f

# docker compose down
down:
	docker compose down --rmi local

# rm  binary		
clean:
	echo "remove bin exe"
	rm -f $(addprefix $(BIN_DIR)/, $(APP_NAME))