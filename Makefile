# 'strip_debug' will increase the Performance of kubenab
# You can disable this by setting `STRIP_DEBUG=""`
STRIP_DEBUG ?= "-tags 'strip_debug'"

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
	cd cmd/kubenab && go build $(STRIP_DEBUG) -a --installsuffix cgo --ldflags="-s" -o ../../bin/kubenab
