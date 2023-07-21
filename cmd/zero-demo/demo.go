// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// demo is a standard, specification-compliant demo example of the zero service.
package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/superproj/zero/cmd/zero-demo/app"
)

func main() {
	app.NewApp("zero-demo").Run()
}
