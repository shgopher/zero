#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


# This script generates all go files from the corresponding protobuf files.
# Usage: `scripts/update-generated-protobuf.sh`.

set -o errexit
set -o nounset
set -o pipefail

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..

# NOTE: All output from this script needs to be copied back to the calling
# source tree.  This is managed in kube::build::copy_output in build/common.sh.
# If the output set is changed update that function.

APIROOTS=${APIROOTS:-$(git grep --files-with-matches -e '// +k8s:protobuf-gen=package' cmd pkg | \
	xargs -n 1 dirname | \
	sed 's,^,github.com/superproj/zero/,;s,k8s.io/kubernetes/staging/src/,,' | \
	sort | uniq
)}

${ZROOT}/scripts/update-generated-protobuf-dockerized.sh "${APIROOTS}" "$@"

# ex: ts=2 sw=2 et filetype=sh
