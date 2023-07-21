// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package tbb

type AddTransactionRequest struct {
	From    string
	To      string
	Value   int
	FromPWD string
}

type AddTransactionResponse struct{}
