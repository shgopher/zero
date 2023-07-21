#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


# This file is not intended to be run automatically. It is meant to be run
# immediately before exporting docs. We do not want to check these documents in
# by default.

set -o errexit
set -o nounset
set -o pipefail

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${ZROOT}/scripts/lib/init.sh"

zero::golang::setup_env

BINS=(
  gen-docs
  gen-man
  gen-zero-docs
  gen-yaml
)
make build -C "${ZROOT}" BINS="${BINS[*]}"

zero::util::ensure-temp-dir

zero::util::gen-docs "${PROJ_TEMP}"

# remove all of the old docs
zero::util::remove-gen-docs

# Copy fresh docs into the repo.
# the shopt is so that we get docs/.generated_docs from the glob.
shopt -s dotglob
cp -af "${PROJ_TEMP}"/* "${ZROOT}"
shopt -u dotglob
