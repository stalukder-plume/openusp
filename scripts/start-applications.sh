#!/bin/bash

# OpenUSP Applications Startup Script
# Starts local OpenUSP applications that connect to containerized infrastructure
# Note: Applications now use YAML configuration with environment variable fallbacks

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting OpenUSP Local Applications...${NC}"

# Set up default environment variables
echo "🔧 Setting up environment variables..."
source "$SCRIPT_DIR/load-env.sh"

# Check if binaries exist
BINARIES=("openusp-controller" "openusp-apiserver" "openusp-cli" "openusp-cwmpacs")
BIN_DIR="$PROJECT_ROOT/build/bin"

for binary in "${BINARIES[@]}"; do
    if [[ ! -f "$BIN_DIR/$binary" ]]; then
        echo -e "${RED}❌ Error: Binary $binary not found. Please run 'make build-all' first.${NC}"
        exit 1
    fi
done

# Function to start a service in the background
start_service() {
    local service_name="$1"
    local binary_path="$2"
    local log_file="$PROJECT_ROOT/logs/${service_name}.log"
    
    mkdir -p "$PROJECT_ROOT/logs"
    
    echo -e "${YELLOW}📦 Starting $service_name...${NC}"
    nohup "$binary_path" > "$log_file" 2>&1 &
    local pid=$!
    echo "$pid" > "$PROJECT_ROOT/logs/${service_name}.pid"
    echo -e "${GREEN}   ✅ $service_name started (PID: $pid, Log: $log_file)${NC}"
}

# Start services
echo ""
echo -e "${GREEN}📋 Starting OpenUSP Services:${NC}"

# Start Controller first
start_service "controller" "$BIN_DIR/openusp-controller"
sleep 2

# Start API Server
start_service "apiserver" "$BIN_DIR/openusp-apiserver"
sleep 2

# Start CWMP ACS Server
start_service "cwmpacs" "$BIN_DIR/openusp-cwmpacs"

echo ""
echo -e "${GREEN}✅ All services started successfully!${NC}"
echo ""
echo -e "${GREEN}📋 Service Details:${NC}"
echo "   Controller:  gRPC on localhost:9001"
echo "   API Server:  HTTP on localhost:8081, Health: http://localhost:8081/health"
echo "   CWMP ACS:    HTTP on localhost:7547, HTTPS on localhost:7548"
echo "   CLI:         Use $BIN_DIR/openusp-cli (connects to API Server)"
echo ""
echo -e "${GREEN}📁 Log files located in: $PROJECT_ROOT/logs/${NC}"
echo ""
echo -e "${YELLOW}🔧 Useful commands:${NC}"
echo "   Check status:  ps aux | grep -E '(controller|apiserver|cwmpacs)'"
echo "   View logs:     tail -f $PROJECT_ROOT/logs/[service].log"
echo "   API Health:    curl http://localhost:8081/health"
echo "   Stop all:      $SCRIPT_DIR/stop-applications.sh"
echo ""