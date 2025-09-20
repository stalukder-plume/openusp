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

# Build the application from source
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

# Build with explicit cross-compilation settings
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w" -o cmd/cli/cli ./cmd/cli

# Deploy the application binary into a lean image
FROM debian:stable AS build-release-stage

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
COPY --from=build-stage /app/cmd/cli/cli /opt/openusp/bin/cli
COPY --chown=openusp:openusp build/common/entrypoint-common.sh /opt/openusp/scripts/entrypoint-common.sh
COPY --chown=openusp:openusp build/cli/entrypoint.sh /opt/openusp/scripts/entrypoint.sh

# Copy YAML configuration files
COPY --chown=openusp:openusp configs/cli.yaml /opt/openusp/configs/cli.yaml

# Make scripts executable
RUN chmod +x /opt/openusp/scripts/*.sh

# Switch to non-root user
USER openusp

ENTRYPOINT ["/opt/openusp/scripts/entrypoint.sh"]
CMD []