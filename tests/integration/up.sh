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

# Run setup SQL script using PG_DSN
psql "$PG_DSN" -f setup.sql

if [ $? -eq 0 ]; then
    echo "Database setup completed successfully"
else
    echo "Error: Database setup failed"
    exit 1
fi

