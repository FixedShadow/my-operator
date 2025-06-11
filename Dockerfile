ARG ARCH=amd64
ARG OS=linux
ARG GOLANG_BUILDER=1.24

FROM quay.io/prometheus/golang-builder:${GOLANG_BUILDER}-base AS builder
WORKDIR /workspace

ENV GOPROXY=https://goproxy.cn,direct

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go mod download -x && go mod verify

ARG GOARGH
ENV GOARGH=${GOARGH}
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build make operator

FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

COPY --from=builder workspace/operator /bin/operator

USER 65534

ENTRYPOINT ["/bin/operator"]

