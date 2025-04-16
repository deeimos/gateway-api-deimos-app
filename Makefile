APP_NAME=gateway-api
CONFIG=./config/local.yaml
MAIN=./cmd/gateway-api/main.go

.PHONY: run build clean migrate

run:
	CONFIG_PATH=$(CONFIG) go run $(MAIN)

build:
	go build -o $(APP_NAME) $(MAIN)

clean:
	rm -f $(APP_NAME)