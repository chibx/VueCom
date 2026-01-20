MIGRATION_DIR = "./cmd/migrator"
SERVER_DIR = "./cmd/server"

.PHONY: do-migrate
migrate:
	@go run $(MIGRATION_DIR)

.PHONY: build-server
build:
	@go build -o bin/server $(SERVER_DIR)
prod:
	@go build -ldflags "-w -s" -o bin/server $(SERVER_DIR)
