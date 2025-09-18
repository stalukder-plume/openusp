#!/bin/bash
# OpenUSP Controller Entrypoint
# Copyright 2023 N4-Networks.com

set -e

# Source the common template
source /opt/openusp/scripts/entrypoint-common.sh

# Service-specific configuration
SERVICE_NAME="Controller"
SERVICE_BINARY="/opt/openusp/bin/controller"
DEFAULT_PORT="8082"
DEFAULT_LOG_LEVEL="info"

# Controller specific environment variables
CONTROLLER_GRPC_PORT="${CONTROLLER_GRPC_PORT:-8083}"
CONTROLLER_MQTT_BROKER="${CONTROLLER_MQTT_BROKER:-mqtt://localhost:1883}"
CONTROLLER_STOMP_BROKER="${CONTROLLER_STOMP_BROKER:-stomp://localhost:61613}"
CONTROLLER_COAP_PORT="${CONTROLLER_COAP_PORT:-5683}"
CONTROLLER_WEBSOCKET_PORT="${CONTROLLER_WEBSOCKET_PORT:-8084}"

# Controller specific help
show_service_help() {
    cat << EOF
Controller specific options:
    --grpc-port PORT      gRPC server port (default: 8083)
    --mqtt-broker URL     MQTT broker URL (default: mqtt://localhost:1883)
    --stomp-broker URL    STOMP broker URL (default: stomp://localhost:61613)
    --coap-port PORT      CoAP server port (default: 5683)
    --ws-port PORT        WebSocket server port (default: 8084)

Controller specific environment variables:
    CONTROLLER_GRPC_PORT     gRPC server port (default: 8083)
    CONTROLLER_MQTT_BROKER   MQTT broker URL
    CONTROLLER_STOMP_BROKER  STOMP broker URL  
    CONTROLLER_COAP_PORT     CoAP server port (default: 5683)
    CONTROLLER_WEBSOCKET_PORT WebSocket server port (default: 8084)
EOF
}

# Controller specific argument parsing
parse_service_args() {
    local remaining_args=()
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --grpc-port)
                CONTROLLER_GRPC_PORT="$2"
                shift 2
                ;;
            --mqtt-broker)
                CONTROLLER_MQTT_BROKER="$2"
                shift 2
                ;;
            --stomp-broker)
                CONTROLLER_STOMP_BROKER="$2"
                shift 2
                ;;
            --coap-port)
                CONTROLLER_COAP_PORT="$2"
                shift 2
                ;;
            --ws-port)
                CONTROLLER_WEBSOCKET_PORT="$2"
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

# Controller specific environment setup
setup_service_environment() {
    export CONTROLLER_GRPC_PORT
    export CONTROLLER_MQTT_BROKER
    export CONTROLLER_STOMP_BROKER
    export CONTROLLER_COAP_PORT
    export CONTROLLER_WEBSOCKET_PORT
}

# Controller specific startup info
show_service_info() {
    echo "gRPC Port: $CONTROLLER_GRPC_PORT"
    echo "MQTT Broker: $CONTROLLER_MQTT_BROKER"
    echo "STOMP Broker: $CONTROLLER_STOMP_BROKER"
    echo "CoAP Port: $CONTROLLER_COAP_PORT"
    echo "WebSocket Port: $CONTROLLER_WEBSOCKET_PORT"
}

# Execute main function from template
main "$@"