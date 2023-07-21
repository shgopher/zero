// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package server

import (
	"context"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"golang.org/x/text/language"

	"github.com/superproj/zero/internal/gateway/locales"
	authmw "github.com/superproj/zero/internal/gateway/server/middleware/auth"
	"github.com/superproj/zero/internal/pkg/idempotent"
	zerometrics "github.com/superproj/zero/internal/pkg/metrics"
	"github.com/superproj/zero/internal/pkg/middleware/auth"
	i18nmw "github.com/superproj/zero/internal/pkg/middleware/i18n"
	idempotentmw "github.com/superproj/zero/internal/pkg/middleware/idempotent"
	"github.com/superproj/zero/internal/pkg/middleware/logging"
	"github.com/superproj/zero/internal/pkg/middleware/tracing"
	"github.com/superproj/zero/internal/pkg/middleware/validate"
	"github.com/superproj/zero/pkg/i18n"
)

// ProviderSet defines a wire provider set.
var ProviderSet = wire.NewSet(NewServers, NewGRPCServer, NewHTTPServer, NewMiddlewares)

// NewServers is a wire provider function that creates and returns a slice of transport servers.
func NewServers(hs *http.Server, gs *grpc.Server) []transport.Server {
	return []transport.Server{hs, gs}
}

func NewWhiteListMatcher() selector.MatchFunc {
	whitelist := make(map[string]struct{})
	// Placeholder
	// whitelist[v1.OperationGatewayGetMiner] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whitelist[operation]; ok {
			return false
		}
		return true
	}
}

func NewMiddlewares(logger log.Logger, idt *idempotent.Idempotent, a auth.AuthProvider, v validate.ICustomValidator) []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(
			recovery.WithLogger(logger),
			recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
				// Handle the panic and return nil.
				return nil
			}),
		),
		metrics.Server(
			metrics.WithSeconds(prom.NewHistogram(zerometrics.KratosMetricSeconds)),
			metrics.WithRequests(prom.NewCounter(zerometrics.KratosServerMetricRequests)),
		),
		i18nmw.Translator(i18n.WithLanguage(language.English), i18n.WithFS(locales.Locales)),
		// circuitbreaker.Client(),
		idempotentmw.Idempotent(idt),
		ratelimit.Server(),
		tracing.Server(),
		selector.Server(authmw.Auth(a)).Match(NewWhiteListMatcher()).Build(),
		validate.Validator(v),
		logging.Server(logger),
	}
}
