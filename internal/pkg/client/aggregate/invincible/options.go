// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package invincible

import (
	"github.com/spf13/pflag"

	gwstore "github.com/superproj/zero/internal/gateway/store"
	"github.com/superproj/zero/internal/pkg/client/usercenter"
	ucstore "github.com/superproj/zero/internal/usercenter/store"
)

// InvincibleOptions is a list of options for all available
// interfaces on the zero platform.
type InvincibleOptions struct {
	UserCenterOptions *usercenter.UserCenterOptions
	GatewayStore      gwstore.IStore
	UserCenterStore   ucstore.IStore
}

// NewInvincibleOptions returns initialized InvincibleOptions.
func NewInvincibleOptions() *InvincibleOptions {
	return &InvincibleOptions{
		UserCenterOptions: usercenter.NewUserCenterOptions(),
	}
}

// Validate validates all the required options.
func (o *InvincibleOptions) Validate() []error {
	if o == nil {
		return nil
	}

	allErrs := []error{}
	allErrs = append(allErrs, o.UserCenterOptions.Validate()...)

	return allErrs
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet.
func (o *InvincibleOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	o.UserCenterOptions.AddFlags(fs)
}
