# 'strip_debug' will increase the Performance of kubenab
# You can disable this by setting `STRIP_DEBUG=""`
STRIP_DEBUG ?=-tags 'strip_debug'
LDFLAGS := "-X main.version=${VERSION}"
# OUT_DIR sets the Path where the kubenab Build Artifact will be puttet
OUT_DIR ?=../../bin
GIT_HASH=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
APP_VERSION=$(shell git describe --abbrev=0 --tags)

.PHONY: image
image:
	@echo "++ Building kubenab docker image..."
	docker build -t kubenab .

.PHONY: build
build: export GOARCH=amd64
build: export CGO_ENABLED=0
build: export GO111MODULE=on
build: export GOPROXY=https://gocenter.io
build:
	@echo "++ Building kubenab go binary..."
	mkdir -p bin
	cd cmd/kubenab && go mod download && \
		go build $(STRIP_DEBUG) -ldflags $(LDFLAGS) -o $(OUT_DIR)/kubenab
