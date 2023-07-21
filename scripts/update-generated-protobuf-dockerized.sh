#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


# This script genertates `*/api.pb.go` from the protobuf file `*/api.proto`.
# Usage: 
#     scripts/update-generated-protobuf-dockerized.sh "${APIROOTS}"
#     An example APIROOT is: "k8s.io/api/admissionregistration/v1"

set -o errexit
set -o nounset
set -o pipefail

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${ZROOT}/scripts/lib/init.sh"

zero::golang::setup_env

GO111MODULE=on go install k8s.io/code-generator/cmd/go-to-protobuf@latest
GO111MODULE=on go install k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo@latest

if [[ -z "$(which protoc)" ]]; then
  echo "Generating protobuf requires protoc 3.0.0-beta1 or newer. Please download and"
  echo "install the platform appropriate Protobuf package for your OS: "
  echo
  echo "  https://github.com/protocolbuffers/protobuf/releases"
  echo
  echo "WARNING: Protobuf changes are not being validated"
  exit 1
fi

gotoprotobuf=go-to-protobuf

while IFS=$'\n' read -r line; do
  APIROOTS+=( "$line" );
done <<< "${1}"
shift

# requires the 'proto' tag to build (will remove when ready)
# searches for the protoc-gen-gogo extension in the output directory
# satisfies import of github.com/gogo/protobuf/gogoproto/gogo.proto and the
# core Google protobuf types
PATH="${ZROOT}/_output/bin:${PATH}" \
  "${gotoprotobuf}" \
  --proto-import="${ZROOT}/third_party" \
  --packages="$(IFS=, ; echo "${APIROOTS[*]}")" \
  --go-header-file "${ZROOT}/scripts/boilerplate.go.txt" \
  "$@"
