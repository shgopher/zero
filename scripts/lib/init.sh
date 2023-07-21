#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


set -o errexit
set +o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
# https://github.com/minerrnetes/minerrnetes/issues/52255
unset CDPATH

# Default use go modules
export GO111MODULE=on

# The root of the build/dist directory
ZROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${ZROOT}/scripts/lib/util.sh"
source "${ZROOT}/scripts/lib/logging.sh"
source "${ZROOT}/scripts/lib/color.sh"

kube::log::install_errexit

source "${ZROOT}/scripts/lib/version.sh"
source "${ZROOT}/scripts/lib/golang.sh"

# list of all available group versions. This should be used when generated code
# or when starting an API server that you want to have everything.
# most preferred version for a group should appear first
PROJ_AVAILABLE_GROUP_VERSIONS="${PROJ_AVAILABLE_GROUP_VERSIONS:-\
apps/v1beta1 \
coordination/v1 \
core/v1 \
}"
