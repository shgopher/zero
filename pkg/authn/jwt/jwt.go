// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package jwt

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt/v4"

	"github.com/superproj/zero/pkg/authn"
)

const (
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"

	// defaultKey holds the default key used to sign a jwt token.
	defaultKey = "zero666"
)

var (
	ErrTokenInvalid           = errors.Unauthorized(reason, "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized(reason, "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized(reason, "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized(reason, "Wrong signing method")
	ErrSignToken              = errors.Unauthorized(reason, "Can not sign token. Is the key correct?")
)

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       2 * time.Hour,
	signingMethod: jwt.SigningMethodHS256,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(defaultKey), nil
	},
}

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       time.Duration
	tokenType     string
	tokenHeader   map[string]interface{}
}

// Option is jwt option.
type Option func(*options)

// WithSigningMethod set signature method.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithSigningKey set the signature key.
func WithSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// WithKeyfunc set the callback function for verifying the key.
func WithKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// WithExpired set the token expiration time (in seconds, default 2h).
func WithExpired(expired time.Duration) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// WithTokenHeader set the customer tokenHeader for client side.
func WithTokenHeader(header map[string]interface{}) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

// New create a authentication instance.
func New(store Storer, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &JWTAuth{opts: &o, store: store}
}

// JWTAuth implement the authn.Authenticator interface.
type JWTAuth struct {
	opts  *options
	store Storer
}

// Sign is used to generate a token.
func (a *JWTAuth) Sign(ctx context.Context, userID string) (authn.IToken, error) {
	now := time.Now()
	expiresAt := now.Add(a.opts.expired)

	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(now),
		Subject:   userID,
	})
	if a.opts.tokenHeader != nil {
		for k, v := range a.opts.tokenHeader {
			token.Header[k] = v
		}
	}

	accessToken, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, ErrSignToken
	}

	tokenInfo := &tokenInfo{
		ExpiresAt: expiresAt.Unix(),
		Type:      a.opts.tokenType,
		Token:     accessToken,
	}

	return tokenInfo, nil
}

// parseToken is used to parse the input accessToken.
func (a *JWTAuth) parseToken(accessToken string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, a.opts.keyfunc)
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if !ok {
			return nil, errors.Unauthorized(reason, err.Error())
		}
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, ErrTokenInvalid
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenParseFail
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	if token.Method != a.opts.signingMethod {
		return nil, ErrUnSupportSigningMethod
	}

	return token.Claims.(*jwt.RegisteredClaims), nil
}

func (a *JWTAuth) callStore(fn func(Storer) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

// Destroy is used to destroy a token.
func (a *JWTAuth) Destroy(ctx context.Context, accessToken string) error {
	claims, err := a.parseToken(accessToken)
	if err != nil {
		return err
	}

	// If storage is set, put the unexpired token in
	store := func(store Storer) error {
		expired := time.Until(claims.ExpiresAt.Time)
		return store.Set(ctx, accessToken, expired)
	}
	return a.callStore(store)
}

// ParseClaims parse the token and return the claims.
func (a *JWTAuth) ParseClaims(ctx context.Context, accessToken string) (*jwt.RegisteredClaims, error) {
	if accessToken == "" {
		return nil, ErrTokenInvalid
	}

	claims, err := a.parseToken(accessToken)
	if err != nil {
		return nil, err
	}

	store := func(store Storer) error {
		exists, err := store.Check(ctx, accessToken)
		if err != nil {
			return err
		}

		if exists {
			return ErrTokenInvalid
		}

		return nil
	}

	if err := a.callStore(store); err != nil {
		return nil, err
	}

	return claims, nil
}

// Release used to release the requested resources.
func (a *JWTAuth) Release() error {
	return a.callStore(func(store Storer) error {
		return store.Close()
	})
}
