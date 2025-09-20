#!/bin/bash

# OpenUSP Infrastructure Startup Script
# Starts containerized MongoDB, ActiveMQ, and Redis services for local development

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$PROJECT_ROOT/deployments/docker-compose-local-dev.yaml"

echo "🚀 Starting OpenUSP Infrastructure Services..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Error: Docker is not running. Please start Docker first."
    exit 1
fi

# Check if compose file exists
if [[ ! -f "$COMPOSE_FILE" ]]; then
    echo "❌ Error: Docker Compose file not found at $COMPOSE_FILE"
    exit 1
fi

# Start infrastructure services
echo "📦 Starting MongoDB, ActiveMQ, Redis, and Swagger UI containers..."
cd "$PROJECT_ROOT"
export PROJECT_ROOT
docker-compose -f "$COMPOSE_FILE" up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 5

# Check service status
echo "🔍 Checking service status..."
docker-compose -f "$COMPOSE_FILE" ps

# Display connection information
echo ""
echo "✅ Infrastructure services are starting up!"
echo ""
echo "📋 Service Connection Details:"
echo "   MongoDB:   localhost:27017 (admin/admin)"
echo "   ActiveMQ:  localhost:61613 (STOMP), localhost:1883 (MQTT)"
echo "   Redis:     localhost:6379"
echo "   Swagger UI: http://localhost:8080"
echo "   ActiveMQ Console: http://localhost:8161/admin (admin/admin)"
echo ""
echo "🔧 To use these services with local OpenUSP apps:"
echo "   source $PROJECT_ROOT/scripts/load-env.sh"
echo ""
echo "🛑 To stop infrastructure services:"
echo "   docker-compose -f $COMPOSE_FILE down"
echo ""