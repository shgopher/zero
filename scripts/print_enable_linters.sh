#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


tmpdir=$(mktemp -d)

function disable_linters() {
  cat << EOF
golint
tagliatelle
wrapcheck
forcetypeassert
goerr113
gomnd
wsl
testpackage
gochecknoglobals
interfacer
maligned
scopelint
gocritic
EOF
}

disable_linters | sort > ${tmpdir}/disable_linters
golangci-lint linters | awk -F':| ' '!match($0, /Enabled|Disabled|^$/){print $1}' | sort > ${tmpdir}/all_linters

for linter in $(comm -3 ${tmpdir}/all_linters ${tmpdir}/disable_linters)
do
  echo "    - $linter"
done

rm -rf ${tmpdir}
