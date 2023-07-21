// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	componentbasevalidation "k8s.io/component-base/config/validation"

	"github.com/superproj/zero/internal/controller/minerset/apis/config"
	"github.com/superproj/zero/internal/pkg/util/validation"
)

// Validate ensures validation of the MinerSetControllerConfiguration struct.
func Validate(cc *config.MinerSetControllerConfiguration) field.ErrorList {
	allErrs := field.ErrorList{}
	newPath := field.NewPath("MinerSetControllerConfiguration")

	effectiveFeatures := utilfeature.DefaultFeatureGate.DeepCopy()
	if err := effectiveFeatures.SetFromMap(cc.FeatureGates); err != nil {
		allErrs = append(allErrs, field.Invalid(newPath.Child("featureGates"), cc.FeatureGates, err.Error()))
	}
	allErrs = append(allErrs, componentbasevalidation.ValidateLeaderElectionConfiguration(&cc.LeaderElection, field.NewPath("leaderElection"))...)

	/*
			if config.ConfigSyncPeriod.Duration <= 0 {
			        allErrs = append(allErrs, field.Invalid(newPath.Child("ConfigSyncPeriod"), config.ConfigSyncPeriod, "must be greater than 0"))
				    }

		if netutils.ParseIPSloppy(cc.BindAddress) == nil {
			allErrs = append(allErrs, field.Invalid(newPath.Child("BindAddress"), config.BindAddress, "not a valid textual representation of an IP address"))
		}
	*/

	if cc.HealthzBindAddress != "" {
		allErrs = append(allErrs, validation.ValidateHostPort(cc.HealthzBindAddress, newPath.Child("healthzBindAddress"))...)
	}

	if cc.MetricsBindAddress != "" {
		allErrs = append(allErrs, validation.ValidateHostPort(cc.MetricsBindAddress, newPath.Child("metricsBindAddress"))...)
	}

	return allErrs
}
