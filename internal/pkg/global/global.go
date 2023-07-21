// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package global

import "sync"

var mux sync.Mutex

func SetDB() {
	mux.Lock()
	defer mux.Unlock()
}
