#!/bin/bash
# OpenUSP OB-USP-A Entrypoint
# Copyright 2023 N4-Networks.com

set -e

# Source the common template
source /opt/obuspa/scripts/entrypoint-common.sh

# Service-specific configuration
SERVICE_NAME="OB-USP-A"
SERVICE_BINARY="/opt/obuspa/bin/obuspa"
DEFAULT_PORT="8080"
DEFAULT_LOG_LEVEL="info"

# OB-USP-A specific environment variables
OBUSPA_HTTP_PORT="${OBUSPA_HTTP_PORT:-8080}"
OBUSPA_HTTPS_PORT="${OBUSPA_HTTPS_PORT:-8443}"
OBUSPA_DATABASE_FILE="${OBUSPA_DATABASE_FILE:-/var/lib/obuspa/usp.db}"
OBUSPA_FACTORY_RESET_FILE="${OBUSPA_FACTORY_RESET_FILE:-}"
OBUSPA_TRUST_STORE_DIR="${OBUSPA_TRUST_STORE_DIR:-/opt/obuspa/configs/certs}"

# OB-USP-A specific help
show_service_help() {
    cat << EOF
OB-USP-A specific options:
    --http-port PORT      HTTP server port (default: 8080)
    --https-port PORT     HTTPS server port (default: 8443)
    --database FILE       Database file path (default: /var/lib/obuspa/usp.db)
    --factory-reset FILE  Factory reset configuration file
    --trust-store DIR     Trust store directory for certificates

OB-USP-A specific environment variables:
    OBUSPA_HTTP_PORT         HTTP server port (default: 8080)
    OBUSPA_HTTPS_PORT        HTTPS server port (default: 8443)
    OBUSPA_DATABASE_FILE     Database file path
    OBUSPA_FACTORY_RESET_FILE Factory reset configuration file
    OBUSPA_TRUST_STORE_DIR   Trust store directory
EOF
}

# OB-USP-A specific argument parsing
parse_service_args() {
    local remaining_args=()
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --http-port)
                OBUSPA_HTTP_PORT="$2"
                shift 2
                ;;
            --https-port)
                OBUSPA_HTTPS_PORT="$2"
                shift 2
                ;;
            --database)
                OBUSPA_DATABASE_FILE="$2"
                shift 2
                ;;
            --factory-reset)
                OBUSPA_FACTORY_RESET_FILE="$2"
                shift 2
                ;;
            --trust-store)
                OBUSPA_TRUST_STORE_DIR="$2"
                shift 2
                ;;
            *)
                remaining_args+=("$1")
                shift 1
                ;;
        esac
    done
    
    echo "${remaining_args[@]}"
}

# OB-USP-A specific environment setup
setup_service_environment() {
    export OBUSPA_HTTP_PORT
    export OBUSPA_HTTPS_PORT
    export OBUSPA_DATABASE_FILE
    
    if [ -n "$OBUSPA_FACTORY_RESET_FILE" ]; then
        export OBUSPA_FACTORY_RESET_FILE
    fi
    
    if [ -n "$OBUSPA_TRUST_STORE_DIR" ]; then
        export OBUSPA_TRUST_STORE_DIR
    fi
    
    # Create database directory if it doesn't exist
    mkdir -p "$(dirname "$OBUSPA_DATABASE_FILE")"
    
    # Create trust store directory if it doesn't exist
    if [ -n "$OBUSPA_TRUST_STORE_DIR" ]; then
        mkdir -p "$OBUSPA_TRUST_STORE_DIR"
    fi
}

# OB-USP-A specific startup info
show_service_info() {
    echo "HTTP Port: $OBUSPA_HTTP_PORT"
    echo "HTTPS Port: $OBUSPA_HTTPS_PORT"
    echo "Database File: $OBUSPA_DATABASE_FILE"
    echo "Factory Reset File: ${OBUSPA_FACTORY_RESET_FILE:-none}"
    echo "Trust Store Directory: ${OBUSPA_TRUST_STORE_DIR:-none}"
}

# Execute main function from template
main "$@"