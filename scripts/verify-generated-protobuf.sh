#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


# This script checks whether updating of the packages which explicitly
# requesting generation is needed or not. We should run
# `scripts/update-generated-protobuf.sh` if those packages are out of date.
# Usage: `scripts/verify-generated-protobuf.sh`.

set -o errexit
set -o nounset
set -o pipefail

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${ZROOT}/scripts/lib/init.sh"

zero::golang::setup_env

APIROOTS=$({
  # gather the packages explicitly requesting generation
  git grep --files-with-matches -e '// +k8s:protobuf-gen=package' cmd pkg | xargs -n 1 dirname
} | sort | uniq)

_tmp="${ZROOT}/_output/_tmp"

cleanup() {
  rm -rf "${_tmp}"
}

trap "cleanup" EXIT SIGINT

cleanup
for APIROOT in ${APIROOTS}; do
  mkdir -p "${_tmp}/${APIROOT}"
  cp -a "${ZROOT}/${APIROOT}"/* "${_tmp}/${APIROOT}/"
done

KUBE_VERBOSE=3 "${ZROOT}/scripts/update-generated-protobuf.sh"
for APIROOT in ${APIROOTS}; do
  TMP_APIROOT="${_tmp}/${APIROOT}"
  echo "diffing ${APIROOT} against freshly generated protobuf"
  ret=0
  diff -qNaupr -I 'Auto generated by' -x 'zz_generated.*' -x '.github' -x '.import-restrictions' "${ZROOT}/${APIROOT}" "${TMP_APIROOT}" || ret=$?
  cp -a "${TMP_APIROOT}"/* "${ZROOT}/${APIROOT}/"
  if [[ $ret -eq 0 ]]; then
    echo "${APIROOT} up to date."
  else
    #echo "${APIROOT} is out of date. Please run scripts/update-generated-protobuf.sh"
    KUBE_VERBOSE=3 "${ZROOT}/scripts/update-generated-protobuf.sh"
    #exit 1
  fi
done
