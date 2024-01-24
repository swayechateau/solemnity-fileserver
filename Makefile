BINARY=fileserver

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	rm -f ${BINARY}
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build: builds the file server binary as a linux executable
build:
	@echo "Building file server binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${BINARY} ./
	@echo "Done!"
