# Variables
APP_NAME=main
DOCKER_COMPOSE=docker-compose
DOCKER_FILE=Dockerfile

# Go commands
.PHONY: build run fmt lint test clean

build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

run:
	go run ./cmd/$(APP_NAME)

fmt:
	go fmt ./...

lint:
	golangci-lint run || true

test:
	go test ./...

clean:
	rm -rf bin/

# Docker commands
.PHONY: docker-build docker-up docker-down docker-logs

docker-build:
	docker build -t $(APP_NAME):latest -f $(DOCKER_FILE) .

docker-up:
	$(DOCKER_COMPOSE) up -d --build

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Prometheus
.PHONY: prometheus

prometheus:
	@echo "Prometheus running at http://localhost:9090"
	xdg-open http://localhost:9090 || open http://localhost:9090 || true

# Full development workflow
.PHONY: dev

dev: fmt lint test build docker-up prometheus
