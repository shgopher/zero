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

KINDS=(deployment service)

for component in $(ls ${ZROOT}/cmd)
do
  truncate -s 0 ${ZROOT}/deployments/${component}.yaml

  for kind in ${KINDS[@]}
  do
    echo -e "---\n# Source: deployments/${component}-${kind}.yaml" >> ${ZROOT}/deployments/${component}.yaml
    sed '/^#\|^$/d' ${ZROOT}/deployments/${component}-${kind}.yaml >> ${ZROOT}/deployments/${component}.yaml
  done

  zero::log::info "generate ${ZROOT}/deployments/${component}.yaml success"
done
