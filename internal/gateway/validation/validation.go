// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package validation

import (
	"context"

	"github.com/google/wire"

	"github.com/superproj/zero/internal/gateway/store"
	"github.com/superproj/zero/internal/pkg/validation"
	v1 "github.com/superproj/zero/pkg/api/gateway/v1"
)

// ProviderSet is a set of validator providers, used for dependency injection.
var ProviderSet = wire.NewSet(New)

// validator is a struct that implements the validate.ICustomValidator interface.
type validator struct {
	ds store.IStore // Data store instance.
}

// New is a factory function that creates and initializes the custom validator.
// It takes a store.IStore instance as input and returns a validation.CustomValidator instance.
func New(ds store.IStore) (validation.CustomValidator, error) {
	vd := &validator{ds: ds}

	return vd, nil
}

// ValidateListMinerSetRequest is a method that validates the ListMinerSetRequest input.
// In this particular case, no validation is performed and it always returns nil.
func (vd *validator) ValidateListMinerSetRequest(ctx context.Context, req *v1.ListMinerSetRequest) error {
	return nil
}
