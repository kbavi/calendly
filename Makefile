build:
	go build -o bin/rest cmd/rest.go

run: build
	bin/rest

docker-run:
	if [ "$(shell uname)" = "Darwin" ]; then \
		docker-compose up db api -d --build; \
	else \
		docker compose up db api -d --build; \
	fi

docker-down:
	if [ "$(shell uname)" = "Darwin" ]; then \
		docker-compose down; \
	else \
		docker compose down; \
	fi

docker-logs:
	if [ "$(shell uname)" = "Darwin" ]; then \
		docker-compose logs -f; \
	else \
		docker compose logs -f; \
	fi

# New commands for docker tests
docker-test:
	if [ "$(shell uname)" = "Darwin" ]; then \
		docker-compose run --rm test; \
	else \
		docker compose run --rm test; \
	fi

docker-test-build:
	if [ "$(shell uname)" = "Darwin" ]; then \
		docker-compose build test; \
	else \
		docker compose build test; \
	fi

# Run this to build and run tests in one command
docker-integration-tests: docker-test-build docker-test

test-db-setup:
	@echo "Setting up test database..."
	@cd tests/integration && ./up.sh 1>/dev/null

test-db-teardown:
	@echo "Cleaning up test database..."
	@cd tests/integration && ./down.sh 1>/dev/null

integration-tests:
	@make test-db-setup
	@echo "Running integration tests..."
	@go test -v ./tests/integration/...; \
	TEST_EXIT_CODE=$$?; \
	make test-db-teardown; \
	exit $$TEST_EXIT_CODE

.PHONY: integration-tests test-db-setup test-db-teardown docker-run docker-down docker-logs docker-test docker-test-build docker-integration-tests