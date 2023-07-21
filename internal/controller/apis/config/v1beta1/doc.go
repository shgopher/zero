// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// +k8s:deepcopy-gen=package
// +k8s:conversion-gen=github.com/superproj/zero/internal/controller/apis/config
// +k8s:conversion-gen=github.com/superproj/zero/pkg/config/v1beta1
// +k8s:conversion-gen=k8s.io/component-base/config/v1alpha1
// +k8s:defaulter-gen=TypeMeta
// +k8s:defaulter-gen-input=github.com/superproj/zero/internal/controller/apis/config/v1beta1
// +groupName=zerocontrollermanager.config.zero.io

package v1beta1 // import "github.com/superproj/zero/internal/controller/apis/config/v1beta1"
