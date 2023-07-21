// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// The controller manager contains multiple controllers. These controllers monitors
// different resources, and completes different logic.
package main

import (
	"os"

	_ "go.uber.org/automaxprocs"
	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"          // for JSON log format registration
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugin
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration

	"github.com/superproj/zero/cmd/zero-miner-controller/app"
)

func main() {
	command := app.NewControllerCommand()
	code := cli.Run(command)
	os.Exit(code)
}
