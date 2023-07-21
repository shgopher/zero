// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package log

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

type KratosLogger interface {
	// Log implements is used to github.com/go-kratos/kratos/v2/log.Logger interface.
	Log(level log.Level, keyvals ...interface{}) error
}

func (l *zapLogger) Log(level log.Level, keyvals ...interface{}) error {
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		l.z.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	switch level {
	case log.LevelDebug:
		l.z.Sugar().Debugw("", keyvals...)
	case log.LevelInfo:
		l.z.Sugar().Infow("", keyvals...)
	case log.LevelWarn:
		l.z.Sugar().Warnw("", keyvals...)
	case log.LevelError:
		l.z.Sugar().Errorw("", keyvals...)
	case log.LevelFatal:
		l.z.Sugar().Fatalw("", keyvals...)
	}

	return nil
}
