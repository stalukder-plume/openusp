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

APP_NAME = cwmpacs
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')

LDFLAGS = -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

.PHONY: all build clean test

all: build

build:
	@echo "Building $(APP_NAME) version $(VERSION)"
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .

clean:
	@echo "Cleaning up"
	rm -f $(APP_NAME)

test:
	go test ./...

run: build
	./$(APP_NAME)

install:
	go install -ldflags "$(LDFLAGS)" .