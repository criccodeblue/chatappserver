# Load environment variables from .env file
include .env
export

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run: ### Run the app
	go run cmd/chatserver/main.go

build: ### Build the app
	go build -o bin/chatserver cmd/chatserver/main.go

migrate-up: ### Initialise the db
	migrate -path database/migrations/ -database '$(PG_URL)' up

migrate-force-version: ### Force db version
	migrate -path database/migrations -database '$(PG_URL)' force '$(VERSION)'