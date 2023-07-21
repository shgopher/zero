// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package main

import (
	"k8s.io/component-base/cli"

	"github.com/superproj/zero/internal/zeroctl/cmd"
	"github.com/superproj/zero/internal/zeroctl/cmd/util"
)

func main() {
	command := cmd.NewDefaultZeroCtlCommand()
	if err := cli.RunNoErrOutput(command); err != nil {
		// Pretty-print the error and exit with an error.
		util.CheckErr(err)
	}
}
