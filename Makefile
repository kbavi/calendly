build:
	go build -o bin/rest cmd/rest.go

run: build
	bin/rest

docker-run:
	docker-compose up db api -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# New commands for docker tests
docker-test:
	docker-compose run --rm test

docker-test-build:
	docker-compose build test

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