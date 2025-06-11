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

.PHONY: operator
operator:
	go build -o $@ ./cmd/operator/


.PHONY: image
image: GOOS := linux
image: operator-image

.PHONY: operator-image
operator-image:
	$(CONTAINER_CLI) build --build-arg ARCH=$(ARCH) --build-arg GOARCH=$(GOARCH) --build-arg OS=$(GOOS) -t $(IMAGE_OPERATOR):$(TAG) .