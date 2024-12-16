#!/bin/bash

cd ../../
# Load environment variables from .env.test
if [ -f .env ]; then
    export $(cat .env | xargs)
else
    echo "Error: .env file not found"
    exit 1
fi

cd tests/integration

# Check if PG_DSN is set
if [ -z "$PG_DSN" ]; then
    echo "Error: PG_DSN environment variable is not set"
    exit 1
fi

# Drop the test database
psql "$PG_DSN" -c "DROP DATABASE IF EXISTS calendly_test;"

if [ $? -eq 0 ]; then
    echo "Database cleanup completed successfully"
else
    echo "Error: Database cleanup failed"
    exit 1
fi
