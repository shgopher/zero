// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//go:build !windows
// +build !windows

package ospkg

const (
	// PathSeparator is a path separator string for a Unix-like os.
	PathSeparator = "/"

	// NewLine is a new line constant for a Unix-like os.
	NewLine = "\n"
)
