#!/usr/bin/env bash

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${ZROOT}/scripts/lib/init.sh"


if [ $# -ne 2 ];then
    zero::log::error "Usage: gen-dockerfile.sh ${DOCKERFILE_DIR} ${IMAGE_NAME}"
    exit 1
fi

DOCKERFILE_DIR=$1/$2
IMAGE_NAME=$2

declare -A envs

function cat_dockerfile() 
{
	cat << 'EOF'
# syntax=docker/dockerfile:1.4

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for 
# this file is https://github.com/superproj/zero.

# Dockerfile generated by scripts/gen-dockerfile.sh. DO NOT EDIT.

# Build the IMAGE_NAME binary
# Run this with docker build --build-arg prod_image=<golang:x.y.z>
# Default <prod_image> is BASE_IMAGE
ARG prod_image=BASE_IMAGE

FROM ${prod_image}
LABEL maintainer="<colin404@foxmail.com>"

WORKDIR /opt/zero

# Note: the <prod_image> is required to support 
# setting timezone otherwise the build will fail
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY IMAGE_NAME /opt/zero/bin/

ENTRYPOINT ["/opt/zero/bin/IMAGE_NAME"]
EOF
}

function cat_multistage_dockerfile() 
{
	cat << 'EOF'
# syntax=docker/dockerfile:1.4

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for 
# this file is https://github.com/superproj/zero.

# Dockerfile generated by scripts/gen-dockerfile.sh. DO NOT EDIT.

# Build the IMAGE_NAME binary
 
# Production image
# Run this with docker build --build-arg prod_image=<golang:x.y.z>
# Default <prod_image> is BASE_IMAGE
ARG prod_image=BASE_IMAGE

# Ignore Hadolint rule "Always tag the version of an image explicitly."
# It's an invalid finding since the image is explicitly set in the Makefile.
# https://github.com/hadolint/hadolint/wiki/DL3006
# hadolint ignore=DL3006
FROM golang:1.20 as builder
WORKDIR /workspace

# Run this with docker build --build-arg goproxy=$(go env GOPROXY) to override the goproxy
ARG goproxy=https://proxy.golang.org
ARG OS
ARG ARCH

# Run this with docker build .
ENV GOPROXY=$goproxy

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum


# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the sources
COPY api/ api/
COPY cmd/IMAGE_NAME cmd/IMAGE_NAME
COPY pkg/ pkg/
COPY internal/ internal/
COPY third_party/ third_party/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${OS:-linux} GOARCH=${ARCH} go build -a ./cmd/IMAGE_NAME

FROM ${prod_image}
LABEL maintainer="<colin404@foxmail.com>"

WORKDIR /opt/zero

# Note: the <prod_image> is required to support 
# setting timezone otherwise the build will fail
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /workspace/IMAGE_NAME /opt/zero/bin/

ENTRYPOINT ["/opt/zero/bin/IMAGE_NAME"]
EOF
}

function get_base_image() {
  declare -A map=(
    ["zero-fake-miner"]="alpine:3.17"
    ["zeroctl"]="alpine:3.17"
  )

  baseImage=${map[$1]}
  echo ${baseImage:-centos:centos8}
}

if [ ! -d ${DOCKERFILE_DIR} ];then
  mkdir -p ${DOCKERFILE_DIR}
fi

BASE_IMAGE=$(get_base_image ${IMAGE_NAME})

# generate dockerfile
cat_dockerfile | sed -e "s/BASE_IMAGE/${BASE_IMAGE}/g" -e "s/IMAGE_NAME/${IMAGE_NAME}/g" > ${DOCKERFILE_DIR}/Dockerfile

# generate multi-stage dockerfile
cat_multistage_dockerfile | sed -e "s/BASE_IMAGE/${BASE_IMAGE}/g" -e "s/IMAGE_NAME/${IMAGE_NAME}/g" > ${DOCKERFILE_DIR}/Dockerfile.multistage
