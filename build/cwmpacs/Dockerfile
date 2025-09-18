# syntax=docker/dockerfile:1
# Copyright 2023 N4-Networks.com
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the CWMP ACS application from source
FROM --platform=$BUILDPLATFORM golang:1.21 AS build-stage

ARG TARGETOS TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY cmd ./cmd
COPY Makefile ./

# Build CWMP ACS with explicit cross-compilation settings
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w" -o cmd/cwmpacs/cwmpacs ./cmd/cwmpacs

# Deploy the application binary into a lean image
FROM debian:stable AS build-release-stage

# Install minimal runtime dependencies (ca-certificates for HTTPS)
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r openusp && useradd --no-log-init -r -g openusp openusp

# Create application directories following FHS conventions
RUN mkdir -p /opt/openusp/bin \
    /opt/openusp/configs \
    /opt/openusp/scripts \
    /var/lib/openusp \
    /var/log/openusp && \
    chown -R openusp:openusp /opt/openusp /var/lib/openusp /var/log/openusp

WORKDIR /opt/openusp

# Copy the binary and scripts to standard location
COPY --from=build-stage /app/cmd/cwmpacs/cwmpacs /opt/openusp/bin/cwmpacs
COPY --chown=openusp:openusp build/common/entrypoint-common.sh /opt/openusp/scripts/entrypoint-common.sh
COPY --chown=openusp:openusp build/cwmpacs/entrypoint.sh /opt/openusp/scripts/entrypoint.sh

# Make scripts executable
RUN chmod +x /opt/openusp/scripts/*.sh

# Switch to non-root user
USER openusp

# CWMP ACS typically uses ports 7547 (HTTP) and 7548 (HTTPS) for TR-069
EXPOSE 7547 7548

ENTRYPOINT ["/opt/openusp/scripts/entrypoint.sh"]
CMD []