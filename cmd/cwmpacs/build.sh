#!/bin/bash
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

# Build script for CWMP ACS
APP_NAME=cwmpacs
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

echo "Building ${APP_NAME} version ${VERSION} at ${BUILD_TIME}"

go build -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" -o ${APP_NAME} .

echo "Build complete: ${APP_NAME}"