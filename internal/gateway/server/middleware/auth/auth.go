// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package auth

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/superproj/zero/internal/gateway/locales"
	"github.com/superproj/zero/internal/pkg/middleware/auth"
	jwtutil "github.com/superproj/zero/internal/pkg/util/jwt"
	"github.com/superproj/zero/internal/pkg/zerox"
	"github.com/superproj/zero/pkg/api/zerrors"
	"github.com/superproj/zero/pkg/i18n"
	"github.com/superproj/zero/pkg/log"
)

// Auth is a authentication and authorization middleware.
func Auth(a auth.AuthProvider) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			accessToken := jwtutil.TokenFromServerContext(ctx)
			if tr, ok := transport.FromServerContext(ctx); ok {
				userID, allowed, err := a.Auth(ctx, accessToken, "*", tr.Operation())
				if err != nil {
					log.Errorw(err, "Authorization failure occurs", "operation", tr.Operation())
					return nil, err
				}
				if !allowed {
					return nil, zerrors.ErrorForbidden(i18n.FromContext(ctx).T(locales.NoPermission))
				}
				ctx = zerox.NewUserID(ctx, userID)
				ctx = zerox.NewAccessToken(ctx, accessToken)
				ctx = log.WithContext(ctx, "user.id", userID)
			}

			return handler(ctx, req)
		}
	}
}
