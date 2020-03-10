# 'strip_debug' will increase the Performance of kubenab
# You can disable this by setting `STRIP_DEBUG=""` – which is not recommended.
STRIP_DEBUG ?=-tags 'strip_debug'
# OUT_DIR sets the Path where the kubenab Build Artifact will be puttet
OUT_DIR ?=../../../bin
GIT_HASH=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
APP_VERSION=$(shell git describe --abbrev=0 --tags)

.PHONY: image
image:
	@echo "++ Building kubenab docker image..."
	docker build -t kubenab .

.PHONY: build
build:
	@export GOARCH=amd64
	@export CGO_ENABLED=0
	@export GO111MODULE=on
	@export GOPROXY=https://gocenter.io

	@echo "++ Pulling Git Tags ++"
	git fetch --tags
	@echo "++ Building kubenab go binary..."
	mkdir -p bin
	go build $(STRIP_DEBUG) \
		-a --installsuffix cgo \
		--ldflags="-s -X main.version=$(APP_VERSION) -X main.date=$(BUILD_DATE) -X main.commit=$(GIT_HASH)" \
		-o $(OUT_DIR)/kubenab
