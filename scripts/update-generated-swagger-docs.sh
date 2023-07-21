#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


# This script generates `types_swagger_doc_generated.go` files for API group
# versions. That file contains functions on API structs that return
# the comments that should be surfaced for the corresponding API type
# in our API docs.

set -o errexit
set -o nounset
set -o pipefail

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${ZROOT}/scripts/lib/init.sh"
source "${ZROOT}/scripts/lib/swagger.sh"

zero::golang::setup_env

IFS=" " read -r -a GROUP_VERSIONS <<< "${PROJ_AVAILABLE_GROUP_VERSIONS}"

# To avoid compile errors, remove the currently existing files.
for group_version in "${GROUP_VERSIONS[@]}"; do
  rm -f "$(kube::util::group-version-to-pkg-path "${group_version}")/types_swagger_doc_generated.go"
done
# ensure we have the latest genswaggertypedocs built
go install k8s.io/kubernetes/cmd/genswaggertypedocs
for group_version in "${GROUP_VERSIONS[@]}"; do
  kube::swagger::gen_types_swagger_doc "${group_version}" "$(kube::util::group-version-to-pkg-path "${group_version}")"
done
