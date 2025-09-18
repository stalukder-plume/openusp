#!/bin/bash
# OpenUSP CWMP ACS Entrypoint
# Copyright 2023 N4-Networks.com

set -e

# Source the common template
source /opt/openusp/scripts/entrypoint-common.sh

# Service-specific configuration
SERVICE_NAME="CWMP ACS"
SERVICE_BINARY="/opt/openusp/bin/cwmpacs"
DEFAULT_PORT="7547"
DEFAULT_LOG_LEVEL="info"

# CWMP ACS specific environment variables
CWMP_ACS_HTTP_PORT="${CWMP_ACS_HTTP_PORT:-7547}"
CWMP_ACS_HTTPS_PORT="${CWMP_ACS_HTTPS_PORT:-7548}"
CWMP_ACS_TLS_ENABLED="${CWMP_ACS_TLS_ENABLED:-false}"
CWMP_ACS_CERT_FILE="${CWMP_ACS_CERT_FILE:-}"
CWMP_ACS_KEY_FILE="${CWMP_ACS_KEY_FILE:-}"

# CWMP ACS specific help
show_service_help() {
    cat << EOF
CWMP ACS specific options:
    --http-port PORT      HTTP port for ACS (default: 7547)
    --https-port PORT     HTTPS port for ACS (default: 7548)
    --tls-enabled         Enable TLS/HTTPS
    --cert-file FILE      TLS certificate file path
    --key-file FILE       TLS private key file path

CWMP ACS specific environment variables:
    CWMP_ACS_HTTP_PORT   HTTP port (default: 7547)
    CWMP_ACS_HTTPS_PORT  HTTPS port (default: 7548)
    CWMP_ACS_TLS_ENABLED TLS enabled flag (default: false)
    CWMP_ACS_CERT_FILE   TLS certificate file path
    CWMP_ACS_KEY_FILE    TLS private key file path

TR-069 Protocol Information:
    Port 7547 is the standard HTTP port for TR-069 ACS-CPE communication
    Port 7548 is the standard HTTPS port for secure TR-069 communication
EOF
}

# CWMP ACS specific argument parsing
parse_service_args() {
    local remaining_args=()
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --http-port)
                CWMP_ACS_HTTP_PORT="$2"
                shift 2
                ;;
            --https-port)
                CWMP_ACS_HTTPS_PORT="$2"
                shift 2
                ;;
            --tls-enabled)
                CWMP_ACS_TLS_ENABLED="true"
                shift 1
                ;;
            --cert-file)
                CWMP_ACS_CERT_FILE="$2"
                shift 2
                ;;
            --key-file)
                CWMP_ACS_KEY_FILE="$2"
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

# CWMP ACS specific environment setup
setup_service_environment() {
    export CWMP_ACS_HTTP_PORT
    export CWMP_ACS_HTTPS_PORT
    export CWMP_ACS_TLS_ENABLED
    
    if [ -n "$CWMP_ACS_CERT_FILE" ]; then
        export CWMP_ACS_CERT_FILE
    fi
    
    if [ -n "$CWMP_ACS_KEY_FILE" ]; then
        export CWMP_ACS_KEY_FILE
    fi
}

# CWMP ACS specific startup info
show_service_info() {
    echo "HTTP Port: $CWMP_ACS_HTTP_PORT"
    echo "HTTPS Port: $CWMP_ACS_HTTPS_PORT"
    echo "TLS Enabled: $CWMP_ACS_TLS_ENABLED"
    
    if [ "$CWMP_ACS_TLS_ENABLED" = "true" ]; then
        echo "Certificate File: ${CWMP_ACS_CERT_FILE:-none}"
        echo "Key File: ${CWMP_ACS_KEY_FILE:-none}"
    fi
}

# Execute main function from template
main "$@"