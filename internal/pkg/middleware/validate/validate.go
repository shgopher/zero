// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package validate

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/superproj/zero/pkg/api/zerrors"
)

type IValidator interface {
	Validate() error
}

// ICustomValidator defines methods to implement a custom validator.
type ICustomValidator interface {
	Validate(ctx context.Context, req interface{}) error
}

// Validator is a validator middleware.
func Validator(vd ICustomValidator) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(IValidator); ok {
				if err := v.Validate(); err != nil {
					if se := new(errors.Error); errors.As(err, &se) {
						return nil, se
					}

					return nil, zerrors.ErrorInvalidParameter(err.Error()).WithCause(err)
				}
			}

			if err := vd.Validate(ctx, req); err != nil {
				if se := new(errors.Error); errors.As(err, &se) {
					return nil, se
				}

				return nil, zerrors.ErrorInvalidParameter(err.Error()).WithCause(err)
			}

			return handler(ctx, req)
		}
	}
}
