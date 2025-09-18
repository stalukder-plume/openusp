#!/bin/bash
# OpenUSP CLI Entrypoint
# Copyright 2023 N4-Networks.com

set -e

# Source the common template
source /opt/openusp/scripts/entrypoint-common.sh

# Service-specific configuration
SERVICE_NAME="CLI"
SERVICE_BINARY="/opt/openusp/bin/cli"
DEFAULT_PORT="N/A"
DEFAULT_LOG_LEVEL="info"

# CLI specific environment variables
CLI_API_ENDPOINT="${CLI_API_ENDPOINT:-http://localhost:8081}"
CLI_OUTPUT_FORMAT="${CLI_OUTPUT_FORMAT:-table}"
CLI_TIMEOUT="${CLI_TIMEOUT:-30s}"
CLI_INTERACTIVE="${CLI_INTERACTIVE:-false}"

# CLI specific help
show_service_help() {
    cat << EOF
CLI specific options:
    --api-endpoint URL    API server endpoint (default: http://localhost:8081)
    --output FORMAT       Output format: table, json, yaml (default: table)
    --timeout DURATION    Request timeout (default: 30s)
    --interactive         Enable interactive mode
    --non-interactive     Disable interactive mode

CLI specific environment variables:
    CLI_API_ENDPOINT     API server endpoint
    CLI_OUTPUT_FORMAT    Output format (table, json, yaml)
    CLI_TIMEOUT          Request timeout duration
    CLI_INTERACTIVE      Interactive mode enabled

Examples:
    # List devices
    $0 device list

    # Get device info in JSON format
    $0 --output json device get device-123

    # Interactive mode
    $0 --interactive

    # Use different API endpoint
    $0 --api-endpoint http://production-api:8081 device list
EOF
}

# CLI specific argument parsing
parse_service_args() {
    local remaining_args=()
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --api-endpoint)
                CLI_API_ENDPOINT="$2"
                shift 2
                ;;
            --output)
                CLI_OUTPUT_FORMAT="$2"
                shift 2
                ;;
            --timeout)
                CLI_TIMEOUT="$2"
                shift 2
                ;;
            --interactive)
                CLI_INTERACTIVE="true"
                shift 1
                ;;
            --non-interactive)
                CLI_INTERACTIVE="false"
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

# CLI specific environment setup
setup_service_environment() {
    export CLI_API_ENDPOINT
    export CLI_OUTPUT_FORMAT
    export CLI_TIMEOUT
    export CLI_INTERACTIVE
}

# CLI specific startup info
show_service_info() {
    echo "API Endpoint: $CLI_API_ENDPOINT"
    echo "Output Format: $CLI_OUTPUT_FORMAT"
    echo "Timeout: $CLI_TIMEOUT"
    echo "Interactive Mode: $CLI_INTERACTIVE"
}

# Override start_service for CLI (no port info needed)
start_service() {
    echo "Starting $SERVICE_NAME..."
    echo "Log Level: $OPENUSP_LOG_LEVEL"
    
    # Show CLI-specific startup info
    show_service_info
    
    # For CLI, we might want to show help if no arguments provided
    if [ $# -eq 0 ] && [ "$CLI_INTERACTIVE" != "true" ]; then
        echo ""
        echo "No command specified. Use --help for usage information or --interactive for interactive mode."
        echo ""
        exec "$SERVICE_BINARY" --help
    else
        # Start the CLI with arguments
        exec "$SERVICE_BINARY" "$@"
    fi
}

# Execute main function from template
main "$@"