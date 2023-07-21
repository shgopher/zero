#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


ZROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${ZROOT}/scripts/common.sh"

readonly LOCAL_OUTPUT_CONFIGPATH="${LOCAL_OUTPUT_ROOT}/configs"
mkdir -p ${LOCAL_OUTPUT_CONFIGPATH}

cd ${ZROOT}/scripts

export PROJ_APISERVER_INSECURE_BIND_ADDRESS=0.0.0.0
export PROJ_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=0.0.0.0

# 集群内通过kubernetes服务名访问
export PROJ_APISERVER_HOST=miner-apiserver
export PROJ_AUTHZ_SERVER_HOST=miner-authz-server
export PROJ_PUMP_HOST=miner-pump
export PROJ_WATCHER_HOST=miner-watcher

# 配置CA证书路径
export CONFIG_USER_CLIENT_CERTIFICATE=/etc/miner/cert/admin.pem
export CONFIG_USER_CLIENT_KEY=/etc/miner/cert/admin-key.pem
export CONFIG_SERVER_CERTIFICATE_AUTHORITY=/etc/miner/cert/ca.pem

for comp in $(ls ${ZROOT/cmd})
do
  zero::log::info "generate ${LOCAL_OUTPUT_CONFIGPATH}/${comp}.yaml"
  ./gen-config.sh install/environment.sh ../configs/${comp}.yaml > ${LOCAL_OUTPUT_CONFIGPATH}/${comp}.yaml
done

