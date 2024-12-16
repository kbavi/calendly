# Setup Guide

This document provides instructions for setting up, installing, and running the project.

## Prerequisites

- GNU make
- Go 1.22 or higher if running locally
- Docker and Docker Compose if running via docker compose


## Installation

1. Clone the repository:
```sh
git clone https://github.com/kbavi/calendly
cd calendly
```
2. Install dependancies [no required if using docker]
```sh
go mod download
```

## API Documentation
[Postman Collection](./Calendly.postman_collection.json)

## Development Setup

### Docker Setup (Recommended)

This setup uses Docker Compose to run the application and integration tests. The docker compose file will also start a postgres database and run the application.

1. Run integration tests

**This command doesn't stream the test output to the console. The output from the tests will appear on the terminal after the tests are finished.**

```sh
make docker-integration-tests
```

2. Build and run using Docker Compose:
```sh
make docker-build
make docker-up
```
3. The application is now running at http://localhost:8080
4. Shutdown the application
```sh
make docker-down
```

### Local Development

You will need a postgres database connection.

1. Copy the `.env.example` file to `.env` and fill in the `PG_DSN` variable. This lets the application know where to connect to the database.

2. Build the project
```sh
make build
```
3. Run integration tests
```sh
make integration-tests
```
4. Start the application
```sh
make run
```
5. The application is now running at http://localhost:8080
## Additional Resources

- [Project Documentation](./PRD.md)