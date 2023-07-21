#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


# 本脚本功能：根据 scripts/environment.sh 配置，生成 PROJ 组件 YAML 配置文件。
# 示例：gen-config.sh scripts/environment.sh configs/miner-apiserver.yaml

env_file="$1"
template_file="$2"

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${ZROOT}/scripts/lib/init.sh"

if [ $# -ne 2 ];then
    zero::log::error "Usage: gen-config.sh scripts/environment.sh configs/miner-apiserver.yaml"
    exit 1
fi

source "${env_file}"

declare -A envs

set +u
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    if [ -z "$(eval echo \$${env})" ];then
        zero::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" ];then
    zero::log::error 'You may run `source scripts/environment.sh` to set these environment'
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"
