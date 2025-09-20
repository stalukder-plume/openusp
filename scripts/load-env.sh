#!/bin/bash

# Environment setup script for OpenUSP CLI
# Sets default environment variables for local development
# Note: Services now use YAML configuration with environment variable fallbacks

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Setting up OpenUSP environment variables..."

# Database configuration
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-27017}
export DB_USER=${DB_USER:-admin}
export DB_PASSWD=${DB_PASSWD:-admin}
export DB_NAME=${DB_NAME:-usp}

# API Server configuration
export API_SERVER_ADDR=${API_SERVER_ADDR:-http://localhost:8080}
export API_SERVER_AUTH_NAME=${API_SERVER_AUTH_NAME:-admin}
export API_SERVER_AUTH_PASSWD=${API_SERVER_AUTH_PASSWD:-admin}

# Service configuration
export HTTP_PORT=${HTTP_PORT:-8080}
export GRPC_PORT=${GRPC_PORT:-9001}
export ENVIRONMENT=${ENVIRONMENT:-development}

echo "Environment variables set for OpenUSP development"
echo "API Server: $API_SERVER_ADDR"
echo "Database: $DB_HOST:$DB_PORT"
echo "Environment: $ENVIRONMENT"
echo "CLI Auth: $CLI_HTTP_BASIC_AUTH_LOGIN_NAME"