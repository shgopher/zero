// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context, keyvals ...interface{}) context.Context {
	if l := FromContext(ctx); l != nil {
		return l.(*zapLogger).WithContext(ctx, keyvals...)
	}

	return std.WithContext(ctx, keyvals...)
}

func (l *zapLogger) WithContext(ctx context.Context, keyvals ...interface{}) context.Context {
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		return context.WithValue(ctx, contextKey{}, l)
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	return context.WithValue(ctx, contextKey{}, l.With(data...))
}

// C represents for `FromContext` with empty keyvals.
func C(ctx context.Context) Logger {
	return FromContext(ctx)
}

// FromContext returns a logger with predefined values from a context.Context.
func FromContext(ctx context.Context, keyvals ...interface{}) Logger {
	var log Logger = std
	if ctx != nil {
		if logger, ok := ctx.Value(contextKey{}).(Logger); ok {
			log = logger
		}
	}

	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		return log
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	return log.With(data...)
}
