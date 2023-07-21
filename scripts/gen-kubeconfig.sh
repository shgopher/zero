#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


caFile="$1"
certFile="$2"
keyFile="$3"

function gen_kubeconfig() {
  cn=`openssl x509 -in $2 -noout -text|awk -F'CN = ' '/Subject.*CN/{print $NF}'`

	cat << EOF
apiVersion: v1
clusters:
- cluster:
    server: https://127.0.0.1:8443
    certificate-authority-data: `base64 $1 |tr -d '\n'`
  name: ${cn}
contexts:
- context:
    cluster: ${cn}
    user: ${cn}
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: ${cn}
  user:
    client-certificate-data: `base64 $2 |tr -d '\n'`
    client-key-data: `base64 $3 |tr -d '\n'`
EOF
}

gen_kubeconfig ${caFile} ${certFile} ${keyFile}
