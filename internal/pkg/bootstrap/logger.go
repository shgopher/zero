// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package bootstrap

import (
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/superproj/zero/pkg/log"
)

func NewLogger(info AppInfo) klog.Logger {
	return klog.With(log.Default(),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
		"service.id", info.ID,
		"service.name", info.Name,
		"service.version", info.Version,
	)
}
