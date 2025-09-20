#!/bin/bash

# OpenUSP Applications Stop Script
# Stops local OpenUSP applications gracefully

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🛑 Stopping OpenUSP Local Applications...${NC}"

# Services to stop
SERVICES=("controller" "apiserver" "cwmpacs")

# Function to stop a service
stop_service() {
    local service_name="$1"
    local pid_file="$PROJECT_ROOT/logs/${service_name}.pid"
    
    if [[ -f "$pid_file" ]]; then
        local pid=$(cat "$pid_file")
        echo -e "${YELLOW}📦 Stopping $service_name (PID: $pid)...${NC}"
        
        if kill -0 "$pid" 2>/dev/null; then
            # Try graceful shutdown first
            kill -TERM "$pid" 2>/dev/null || true
            sleep 2
            
            # Check if still running
            if kill -0 "$pid" 2>/dev/null; then
                echo -e "${YELLOW}   ⏳ Waiting for graceful shutdown...${NC}"
                sleep 3
                
                # Force kill if still running
                if kill -0 "$pid" 2>/dev/null; then
                    echo -e "${RED}   💥 Force killing $service_name...${NC}"
                    kill -KILL "$pid" 2>/dev/null || true
                fi
            fi
            echo -e "${GREEN}   ✅ $service_name stopped${NC}"
        else
            echo -e "${YELLOW}   ⚠️  $service_name was not running${NC}"
        fi
        
        rm -f "$pid_file"
    else
        echo -e "${YELLOW}📦 $service_name: No PID file found${NC}"
        
        # Try to find and kill by process name
        local pids=$(pgrep -f "openusp-$service_name" || true)
        if [[ -n "$pids" ]]; then
            echo -e "${YELLOW}   🔍 Found running $service_name processes: $pids${NC}"
            echo "$pids" | xargs -r kill -TERM 2>/dev/null || true
            sleep 2
            echo "$pids" | xargs -r kill -KILL 2>/dev/null || true
            echo -e "${GREEN}   ✅ $service_name processes stopped${NC}"
        else
            echo -e "${YELLOW}   ⚠️  No running $service_name processes found${NC}"
        fi
    fi
}

# Stop all services
echo ""
for service in "${SERVICES[@]}"; do
    stop_service "$service"
done

# Clean up any remaining processes
echo ""
echo -e "${YELLOW}🧹 Cleaning up any remaining OpenUSP processes...${NC}"
for service in "${SERVICES[@]}"; do
    pgrep -f "openusp-$service" | xargs -r kill -KILL 2>/dev/null || true
done

echo ""
echo -e "${GREEN}✅ All OpenUSP applications stopped!${NC}"

# Show remaining processes if any
remaining=$(ps aux | grep -E "(openusp-controller|openusp-apiserver|openusp-cwmpacs)" | grep -v grep | grep -v stop-applications || true)
if [[ -n "$remaining" ]]; then
    echo -e "${YELLOW}⚠️  Remaining processes:${NC}"
    echo "$remaining"
else
    echo -e "${GREEN}🎉 All OpenUSP processes have been stopped cleanly${NC}"
fi

echo ""