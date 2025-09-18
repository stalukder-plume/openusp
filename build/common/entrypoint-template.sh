#!/bin/bash
# OpenUSP Common Entrypoint Functions
# Copyright 2023 N4-Networks.com

set -e

# Service-specific defaults (override these in each service's entrypoint.sh)
SERVICE_NAME="${SERVICE_NAME:-openusp-service}"
SERVICE_BINARY="${SERVICE_BINARY:-/opt/openusp/bin/service}"
DEFAULT_PORT="${DEFAULT_PORT:-8080}"
DEFAULT_LOG_LEVEL="${DEFAULT_LOG_LEVEL:-info}"

# Common environment variables with defaults
OPENUSP_PORT="${OPENUSP_PORT:-$DEFAULT_PORT}"
OPENUSP_LOG_LEVEL="${OPENUSP_LOG_LEVEL:-$DEFAULT_LOG_LEVEL}"
OPENUSP_CONFIG_FILE="${OPENUSP_CONFIG_FILE:-}"
OPENUSP_DB_URL="${OPENUSP_DB_URL:-}"

# Function to show standard help
show_help() {
    cat << EOF
OpenUSP $SERVICE_NAME

Usage: $0 [OPTIONS]

Common Options:
    --help                  Show this help message
    --port PORT            Service port (default: $DEFAULT_PORT)
    --log-level LEVEL      Log level: debug, info, warn, error (default: $DEFAULT_LOG_LEVEL)
    --config-file FILE     Configuration file path
    --db-url URL          Database URL override

Common Environment Variables:
    OPENUSP_PORT          Service port (default: $DEFAULT_PORT)
    OPENUSP_LOG_LEVEL     Log level (default: $DEFAULT_LOG_LEVEL)
    OPENUSP_CONFIG_FILE   Configuration file path
    OPENUSP_DB_URL        Database URL

Examples:
    # Start with default settings
    $0

    # Start with custom port
    $0 --port 9090

    # Start with debug logging
    $0 --log-level debug

    # Start with custom config file
    $0 --config-file /opt/openusp/configs/service.yaml

EOF
    # Allow services to add their own help sections
    if declare -f show_service_help > /dev/null; then
        echo "Service-specific options:"
        show_service_help
    fi
}

# Function to parse common arguments (services can extend this)
parse_common_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --help)
                show_help
                exit 0
                ;;
            --port)
                OPENUSP_PORT="$2"
                shift 2
                ;;
            --log-level)
                OPENUSP_LOG_LEVEL="$2"
                shift 2
                ;;
            --config-file)
                OPENUSP_CONFIG_FILE="$2"
                shift 2
                ;;
            --db-url)
                OPENUSP_DB_URL="$2"
                shift 2
                ;;
            *)
                # Return remaining args for service-specific parsing
                echo "$@"
                return
                ;;
        esac
    done
}

# Function to setup common environment
setup_environment() {
    # Set standard environment variables
    export OPENUSP_PORT
    export OPENUSP_LOG_LEVEL
    
    # Set service-specific environment variables if they exist
    if [ -n "$OPENUSP_CONFIG_FILE" ]; then
        export OPENUSP_CONFIG_FILE
    fi
    
    if [ -n "$OPENUSP_DB_URL" ]; then
        export OPENUSP_DB_URL
    fi

    # Create necessary directories
    mkdir -p /var/log/openusp
    mkdir -p /var/lib/openusp
    
    # Allow services to set up their own environment
    if declare -f setup_service_environment > /dev/null; then
        setup_service_environment
    fi
}

# Function to start service (can be overridden)
start_service() {
    echo "Starting $SERVICE_NAME..."
    echo "Port: $OPENUSP_PORT"
    echo "Log Level: $OPENUSP_LOG_LEVEL"
    
    # Allow services to show their own startup info
    if declare -f show_service_info > /dev/null; then
        show_service_info
    fi

    # Start the service
    exec "$SERVICE_BINARY" "$@"
}

# Main execution function (should be called by service entrypoints)
main() {
    # Parse common arguments and get remaining args
    remaining_args=$(parse_common_args "$@")
    
    # Setup environment
    setup_environment
    
    # Allow service-specific argument parsing
    if declare -f parse_service_args > /dev/null; then
        remaining_args=$(parse_service_args $remaining_args)
    fi
    
    # Start the service with remaining arguments
    start_service $remaining_args
}

# Export functions for use by service entrypoints
export -f show_help parse_common_args setup_environment start_service main