// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package auth

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	lru "github.com/hashicorp/golang-lru"
	"gorm.io/gorm"

	known "github.com/superproj/zero/internal/pkg/known/usercenter"
	"github.com/superproj/zero/internal/usercenter/model"
	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
	jwtauthn "github.com/superproj/zero/pkg/authn/jwt"
	"github.com/superproj/zero/pkg/log"
)

const (
	// reasonUnauthorized holds the error reason.
	reasonUnauthorized string = "Unauthorized"
)

var (
	// ErrMissingKID is returned when the token format is invalid and the kid field is missing in the token claims.
	ErrMissingKID = errors.Unauthorized(reasonUnauthorized, "Invalid token format: missing kid field in claims")
	// ErrSecretDisabled is returned when the SecretID is disabled.
	ErrSecretDisabled = errors.Unauthorized(reasonUnauthorized, "SecretID is disabled")
)

// AuthnInterface defines the interface for authentication.
type AuthnInterface interface {
	Verify(accessToken string) (string, error)
}

// AuthnProviderSet is authn providers.
var AuthnProviderSet = wire.NewSet(NewAuthn, wire.Bind(new(AuthnInterface), new(*authn)))

type SecretGetter func(kid string) (*model.SecretM, error)

type authn struct {
	getter  SecretGetter
	secrets *lru.Cache
}

// NewAuthn returns a new instance of authn.
func NewAuthn(getter SecretGetter) (*authn, error) {
	l, err := lru.New(known.DefaultLRUSize)
	if err != nil {
		log.Errorw(err, "Failed to create lru cache")
		return nil, err
	}

	return &authn{getter: getter, secrets: l}, nil
}

// Verify verifies the given access token and returns the userID associated with the token.
func (a *authn) Verify(accessToken string) (string, error) {
	var secret *model.SecretM
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is HMAC signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", jwtauthn.ErrUnSupportSigningMethod
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return "", ErrMissingKID
		}

		var err error
		secret, err = a.GetSecret(kid)
		if err != nil {
			return "", err
		}

		if secret.Status == model.StatusSecretDisabled {
			return "", ErrSecretDisabled
		}

		return []byte(secret.SecretKey), nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if !ok {
			return "", errors.Unauthorized(reasonUnauthorized, err.Error())
		}
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "", jwtauthn.ErrTokenInvalid
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return "", jwtauthn.ErrTokenExpired
		}
		return "", err
	}

	if !token.Valid {
		return "", jwtauthn.ErrTokenInvalid
	}

	if keyExpired(secret.Expires) {
		return "", jwtauthn.ErrTokenExpired
	}

	// you can return claims if you need
	// claims := token.Claims.(*jwt.RegisteredClaims)
	return secret.UserID, nil
}

// GetSecret returns the secret associated with the given key.
func (a *authn) GetSecret(key string) (*model.SecretM, error) {
	s, ok := a.secrets.Get(key)
	if ok {
		return s.(*model.SecretM), nil
	}

	secret, err := a.getter(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrorSecretNotFound(err.Error())
		}

		return nil, err
	}

	a.secrets.Add(key, secret)
	return secret, nil
}
