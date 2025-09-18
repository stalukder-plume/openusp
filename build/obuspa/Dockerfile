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

# Base image with dependencies
FROM debian:stable AS base

# Install dependencies required for building and running obuspa
RUN apt-get update && \
    apt-get install -y \
        build-essential \
        cmake \
        git \
        libcurl4-openssl-dev \
        libmosquitto-dev \
        libprotobuf-c-dev \
        libssl-dev \
        libsqlite3-dev \
        libwebsockets-dev \
        protobuf-c-compiler \
        zlib1g-dev && \
    rm -rf /var/lib/apt/lists/*

# Build stage for obuspa
FROM base AS build-stage

WORKDIR /app

# Copy the obuspa source code
COPY obuspa ./obuspa

# Build obuspa
WORKDIR /app/obuspa
RUN make

# Runtime stage
FROM debian:stable AS runtime

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y \
        libcurl4 \
        libmosquitto1 \
        libprotobuf-c1 \
        libssl3 \
        libsqlite3-0 \
        libwebsockets19 \
        zlib1g && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r obuspa && useradd --no-log-init -r -g obuspa obuspa

# Create application directories following FHS conventions
RUN mkdir -p /opt/obuspa/bin \
    /opt/obuspa/configs \
    /opt/obuspa/scripts \
    /var/lib/obuspa \
    /var/log/obuspa && \
    chown -R obuspa:obuspa /opt/obuspa /var/lib/obuspa /var/log/obuspa

# Copy the obuspa binary and scripts
COPY --from=build-stage /app/obuspa/obuspa /opt/obuspa/bin/obuspa
COPY --chown=obuspa:obuspa build/common/entrypoint-common.sh /opt/obuspa/scripts/entrypoint-common.sh
COPY --chown=obuspa:obuspa build/obuspa/entrypoint.sh /opt/obuspa/scripts/entrypoint.sh

# Make scripts executable
RUN chmod +x /opt/obuspa/scripts/*.sh

WORKDIR /opt/obuspa

# Switch to non-root user
USER obuspa

# Expose default ports (can be overridden)
EXPOSE 8080 8443

# Default entrypoint
ENTRYPOINT ["/opt/obuspa/scripts/entrypoint.sh"]
CMD []