// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package usercenter

import (
	"context"

	"github.com/golang-jwt/jwt/v4"

	"github.com/superproj/zero/internal/usercenter/auth"
	"github.com/superproj/zero/internal/usercenter/model"
	"github.com/superproj/zero/internal/usercenter/store"
	"github.com/superproj/zero/pkg/authn"
	jwtauthn "github.com/superproj/zero/pkg/authn/jwt"
	"github.com/superproj/zero/pkg/authn/jwt/store/redis"
	genericoptions "github.com/superproj/zero/pkg/options"
)

// NewAuthenticator creates a new JWT-based Authenticator using the provided JWT and Redis options.
func NewAuthenticator(jwtOpts *genericoptions.JWTOptions, redisOpts *genericoptions.RedisOptions) (authn.Authenticator, func(), error) {
	// Create a list of options for jwtauthn.
	opts := []jwtauthn.Option{
		jwtauthn.WithExpired(jwtOpts.Expired),
		jwtauthn.WithSigningKey([]byte(jwtOpts.Key)),
		jwtauthn.WithKeyfunc(func(t *jwt.Token) (interface{}, error) {
			// Verify that the signing method is HMAC.
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwtauthn.ErrTokenInvalid
			}
			// Return the signing key.
			return []byte(jwtOpts.Key), nil
		}),
	}

	// Set the signing method based on the provided option.
	var method jwt.SigningMethod
	switch jwtOpts.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}

	opts = append(opts, jwtauthn.WithSigningMethod(method))

	// Create a Redis store for jwtauthn.
	store := redis.NewStore(&redis.Config{
		Addr:      redisOpts.Addr,
		Password:  redisOpts.Password,
		Database:  redisOpts.Database,
		KeyPrefix: "authn_",
	})

	// Create a new jwtauthn instance using the Redis store and options.
	authn := jwtauthn.New(store, opts...)
	// Define a function to release the resources used by jwtauthn.
	// Example of cleanFunc, here we clean nothing.
	cleanFunc := func() {}

	return authn, cleanFunc, nil
}

// newSecretGetter creates a new secret getter function using the provided IStore instance.
func newSecretGetter(ds store.IStore) auth.SecretGetter {
	return func(kid string) (*model.SecretM, error) {
		return ds.Secrets().Get(context.Background(), "", kid)
	}
}
