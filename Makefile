# DEBUG enables or disables if the resulting binary should be compiled with
# debug information data or not.
# Enabling debug will increase resulting file size and decrease performance.
# Default: "false"
DEBUG ?= false

# OUT_DIR sets the Path where the kubenab Build Artifact will be puttet
OUT_DIR ?=./bin
COMMIT := $(shell git rev-parse HEAD)
LD_FLAGS ?=-X internal.Version=$(shell git-semver -prefix v) -X internal.Commit=${COMMIT} -X internal.BuildDate='$(shell date -u +%Y-%m-%d_%T)'
C_FLAGS := ""

# set default target to 'help'
.DEFAULT_GOAL:=help


.PHONY: help
# source: https://blog.thapaliya.com/posts/well-documented-makefiles/
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Building

.PHONY: build
build: ## compile the `kubenab` project
	@export GOARCH=amd64
	@export CGO_ENABLED=0
	@export GO111MODULE=on
	@export GOPROXY=https://gocenter.io

	@echo "++ Pulling Git Tags ++"
	git fetch --tags

	@# strip debug informations if !DEBUG
ifeq "$(DEBUG)" "false"
	strip --strip-debug --strip-unneeded \
		--remove-section='!.go.buildinfo' $(OUT_DIR)/kubenab
else
	@# add 'debug' LD flag
	C_FLAGS+="-tags 'debug'"
endif

	@echo "++ Building kubenab go binary..."
	mkdir -p bin
	go build -a --installsuffix cgo --ldflags="$(LD_FLAGS)" \
		-o $(OUT_DIR)/kubenab



.PHONY: image
image: ## build the Docker Image
	@echo "++ Building kubenab docker image..."
	docker build -t kubenab .

##@ Developing

.PHONY: dev-setup
dev-setup: ## install required tools to get started with developing
	# go one directory backwards to prevent that `git-semver` will be added
	# to `go.(mod|sum)`
	cd ..
	go get github.com/mdomke/git-semver@master
