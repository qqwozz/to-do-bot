.PHONY: help build run clean docker-up docker-down docker-logs bot-run backend-build test

GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

help: ## Show help
	@echo "$(GREEN)To-Do Bot Commands$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

docker-up: ## Start all services
	docker-compose up -d --build
	@echo "$(GREEN)Services started$(NC)"

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show logs
	docker-compose logs -f

backend-build: ## Build C++ backend
	cd cpp-backend && mkdir -p build && cd build && cmake .. && make -j$$(nproc)

backend-run: ## Run backend locally
	cd cpp-backend/build && ./todo-backend

bot-run: ## Run Go bot locally
	go run cmd/bot/main.go

test-go: ## Run Go tests
	go test -v -race ./...

test-cpp: ## Build C++ tests
	cd cpp-backend/tests && mkdir -p build && cd build && cmake .. && make -j$$(nproc)

test-cpp-run: ## Run C++ tests
	cd cpp-backend/tests/build && ./database_tests

test: test-go test-cpp test-cpp-run ## Run all tests

clean: ## Clean build artifacts
	rm -rf cpp-backend/build cpp-backend/tests/build bin *.db
