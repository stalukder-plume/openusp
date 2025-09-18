#!/bin/bash
# OpenUSP API Server Entrypoint
# Copyright 2023 N4-Networks.com

set -e

# Source the common template
source /opt/openusp/scripts/entrypoint-common.sh

# Service-specific configuration
SERVICE_NAME="API Server"
SERVICE_BINARY="/opt/openusp/bin/apiserver"
DEFAULT_PORT="8081"
DEFAULT_LOG_LEVEL="info"

# API Server specific environment variables
APISERVER_HOST="${APISERVER_HOST:-0.0.0.0}"
APISERVER_TLS_ENABLED="${APISERVER_TLS_ENABLED:-false}"
APISERVER_CERT_FILE="${APISERVER_CERT_FILE:-}"
APISERVER_KEY_FILE="${APISERVER_KEY_FILE:-}"
APISERVER_CORS_ENABLED="${APISERVER_CORS_ENABLED:-true}"

# API Server specific help
show_service_help() {
    cat << EOF
API Server specific options:
    --host HOST            Bind host address (default: 0.0.0.0)
    --tls-enabled         Enable TLS/HTTPS
    --cert-file FILE      TLS certificate file path
    --key-file FILE       TLS private key file path
    --cors-enabled        Enable CORS (default: true)

API Server specific environment variables:
    APISERVER_HOST        Bind host address (default: 0.0.0.0)
    APISERVER_TLS_ENABLED TLS enabled flag (default: false)
    APISERVER_CERT_FILE   TLS certificate file path
    APISERVER_KEY_FILE    TLS private key file path
    APISERVER_CORS_ENABLED CORS enabled flag (default: true)
EOF
}

# API Server specific argument parsing
parse_service_args() {
    local remaining_args=()
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --host)
                APISERVER_HOST="$2"
                shift 2
                ;;
            --tls-enabled)
                APISERVER_TLS_ENABLED="true"
                shift 1
                ;;
            --cert-file)
                APISERVER_CERT_FILE="$2"
                shift 2
                ;;
            --key-file)
                APISERVER_KEY_FILE="$2"
                shift 2
                ;;
            --cors-enabled)
                APISERVER_CORS_ENABLED="true"
                shift 1
                ;;
            --cors-disabled)
                APISERVER_CORS_ENABLED="false"
                shift 1
                ;;
            *)
                remaining_args+=("$1")
                shift 1
                ;;
        esac
    done
    
    echo "${remaining_args[@]}"
}

# API Server specific environment setup
setup_service_environment() {
    export APISERVER_HOST
    export APISERVER_TLS_ENABLED
    export APISERVER_CORS_ENABLED
    
    if [ -n "$APISERVER_CERT_FILE" ]; then
        export APISERVER_CERT_FILE
    fi
    
    if [ -n "$APISERVER_KEY_FILE" ]; then
        export APISERVER_KEY_FILE
    fi
}

# API Server specific startup info
show_service_info() {
    echo "Host: $APISERVER_HOST"
    echo "TLS Enabled: $APISERVER_TLS_ENABLED"
    echo "CORS Enabled: $APISERVER_CORS_ENABLED"
    
    if [ "$APISERVER_TLS_ENABLED" = "true" ]; then
        echo "Certificate File: ${APISERVER_CERT_FILE:-none}"
        echo "Key File: ${APISERVER_KEY_FILE:-none}"
    fi
}

# Execute main function from template
main "$@"