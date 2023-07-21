// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package known

const (
	// This exposes compute information based on the miner type.
	CPUAnnotation    = "apps.zero.io/vCPU"
	MemoryAnnotation = "apps.zero.io/memoryMb"
)

const (
	SkipVerifyAnnotation = "apps.zero.io/skip-verify"
)

var AllImmutableAnnotations = []string{
	CPUAnnotation,
	MemoryAnnotation,
}
