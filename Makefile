GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
TAG ?= $(shell git rev-parse --short HEAD)
IMAGE_OPERATOR?=my-operator
ifeq ($(GOARCH),arm)
	ARCH=arm7
else
	ARCH=$(GOARCH)
endif

CONTAINER_CLI ?= docker

GO_BUILD_RECIPE=\
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	CGO_ENABLED=0 \
	go build

.PHONY: operator
operator:
	$(GO_BUILD_RECIPE) -o $@ ./cmd/operator/


.PHONY: image
image: GOOS := linux
image: operator-image

.PHONY: operator-image
operator-image:
	$(CONTAINER_CLI) build --build-arg ARCH=$(ARCH) --build-arg GOARCH=$(GOARCH) --build-arg OS=$(GOOS) -t $(IMAGE_OPERATOR):$(TAG) .