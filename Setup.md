# Setup Guide

This document provides instructions for setting up, installing, and running the project.

## Prerequisites

- Go 1.x or higher
- Docker and Docker Compose
- Make (optional, but recommended)

## Installation

1. Clone the repository:
```sh
git clone https://github.com/kbavi/calendly
cd calendly
```
2. Install dependancies
```sh
go mod download
```

## Development Setup

### Docker Setup (Recommended)

This setup uses Docker Compose to run the application and integration tests. The docker compose file will also start a postgres database and run the application.

1. Run integration tests
```sh
make docker-integration-tests
```

2. Build and run using Docker Compose:
```sh
make docker-build
make docker-up
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
## Additional Resources

- [Project Documentation](./PRD.md)