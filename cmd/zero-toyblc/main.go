// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// zero-toyblc is used to show a naive and simple blockchain.
package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/superproj/zero/cmd/zero-toyblc/app"
)

func main() {
	app.NewApp("zero-toyblc").Run()
}
