//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.

// Code generated by conversion-gen. DO NOT EDIT.

package v1beta1

import (
	unsafe "unsafe"

	config "github.com/superproj/zero/internal/controller/apis/config"
	configv1beta1 "github.com/superproj/zero/pkg/config/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	v1alpha1 "k8s.io/component-base/config/v1alpha1"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*ChainControllerConfiguration)(nil), (*config.ChainControllerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ChainControllerConfiguration_To_config_ChainControllerConfiguration(a.(*ChainControllerConfiguration), b.(*config.ChainControllerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ChainControllerConfiguration)(nil), (*ChainControllerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ChainControllerConfiguration_To_v1beta1_ChainControllerConfiguration(a.(*config.ChainControllerConfiguration), b.(*ChainControllerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*GarbageCollectorControllerConfiguration)(nil), (*config.GarbageCollectorControllerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_GarbageCollectorControllerConfiguration_To_config_GarbageCollectorControllerConfiguration(a.(*GarbageCollectorControllerConfiguration), b.(*config.GarbageCollectorControllerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GarbageCollectorControllerConfiguration)(nil), (*GarbageCollectorControllerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GarbageCollectorControllerConfiguration_To_v1beta1_GarbageCollectorControllerConfiguration(a.(*config.GarbageCollectorControllerConfiguration), b.(*GarbageCollectorControllerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*GenericControllerManagerConfiguration)(nil), (*config.GenericControllerManagerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration(a.(*GenericControllerManagerConfiguration), b.(*config.GenericControllerManagerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GenericControllerManagerConfiguration)(nil), (*GenericControllerManagerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GenericControllerManagerConfiguration_To_v1beta1_GenericControllerManagerConfiguration(a.(*config.GenericControllerManagerConfiguration), b.(*GenericControllerManagerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*GroupResource)(nil), (*config.GroupResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_GroupResource_To_config_GroupResource(a.(*GroupResource), b.(*config.GroupResource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GroupResource)(nil), (*GroupResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GroupResource_To_v1beta1_GroupResource(a.(*config.GroupResource), b.(*GroupResource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ZeroControllerManagerConfiguration)(nil), (*config.ZeroControllerManagerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ZeroControllerManagerConfiguration_To_config_ZeroControllerManagerConfiguration(a.(*ZeroControllerManagerConfiguration), b.(*config.ZeroControllerManagerConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ZeroControllerManagerConfiguration)(nil), (*ZeroControllerManagerConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ZeroControllerManagerConfiguration_To_v1beta1_ZeroControllerManagerConfiguration(a.(*config.ZeroControllerManagerConfiguration), b.(*ZeroControllerManagerConfiguration), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1beta1_ChainControllerConfiguration_To_config_ChainControllerConfiguration(in *ChainControllerConfiguration, out *config.ChainControllerConfiguration, s conversion.Scope) error {
	out.Image = in.Image
	return nil
}

// Convert_v1beta1_ChainControllerConfiguration_To_config_ChainControllerConfiguration is an autogenerated conversion function.
func Convert_v1beta1_ChainControllerConfiguration_To_config_ChainControllerConfiguration(in *ChainControllerConfiguration, out *config.ChainControllerConfiguration, s conversion.Scope) error {
	return autoConvert_v1beta1_ChainControllerConfiguration_To_config_ChainControllerConfiguration(in, out, s)
}

func autoConvert_config_ChainControllerConfiguration_To_v1beta1_ChainControllerConfiguration(in *config.ChainControllerConfiguration, out *ChainControllerConfiguration, s conversion.Scope) error {
	out.Image = in.Image
	return nil
}

// Convert_config_ChainControllerConfiguration_To_v1beta1_ChainControllerConfiguration is an autogenerated conversion function.
func Convert_config_ChainControllerConfiguration_To_v1beta1_ChainControllerConfiguration(in *config.ChainControllerConfiguration, out *ChainControllerConfiguration, s conversion.Scope) error {
	return autoConvert_config_ChainControllerConfiguration_To_v1beta1_ChainControllerConfiguration(in, out, s)
}

func autoConvert_v1beta1_GarbageCollectorControllerConfiguration_To_config_GarbageCollectorControllerConfiguration(in *GarbageCollectorControllerConfiguration, out *config.GarbageCollectorControllerConfiguration, s conversion.Scope) error {
	if err := v1.Convert_Pointer_bool_To_bool(&in.EnableGarbageCollector, &out.EnableGarbageCollector, s); err != nil {
		return err
	}
	out.ConcurrentGCSyncs = in.ConcurrentGCSyncs
	out.GCIgnoredResources = *(*[]config.GroupResource)(unsafe.Pointer(&in.GCIgnoredResources))
	return nil
}

// Convert_v1beta1_GarbageCollectorControllerConfiguration_To_config_GarbageCollectorControllerConfiguration is an autogenerated conversion function.
func Convert_v1beta1_GarbageCollectorControllerConfiguration_To_config_GarbageCollectorControllerConfiguration(in *GarbageCollectorControllerConfiguration, out *config.GarbageCollectorControllerConfiguration, s conversion.Scope) error {
	return autoConvert_v1beta1_GarbageCollectorControllerConfiguration_To_config_GarbageCollectorControllerConfiguration(in, out, s)
}

func autoConvert_config_GarbageCollectorControllerConfiguration_To_v1beta1_GarbageCollectorControllerConfiguration(in *config.GarbageCollectorControllerConfiguration, out *GarbageCollectorControllerConfiguration, s conversion.Scope) error {
	if err := v1.Convert_bool_To_Pointer_bool(&in.EnableGarbageCollector, &out.EnableGarbageCollector, s); err != nil {
		return err
	}
	out.ConcurrentGCSyncs = in.ConcurrentGCSyncs
	out.GCIgnoredResources = *(*[]GroupResource)(unsafe.Pointer(&in.GCIgnoredResources))
	return nil
}

// Convert_config_GarbageCollectorControllerConfiguration_To_v1beta1_GarbageCollectorControllerConfiguration is an autogenerated conversion function.
func Convert_config_GarbageCollectorControllerConfiguration_To_v1beta1_GarbageCollectorControllerConfiguration(in *config.GarbageCollectorControllerConfiguration, out *GarbageCollectorControllerConfiguration, s conversion.Scope) error {
	return autoConvert_config_GarbageCollectorControllerConfiguration_To_v1beta1_GarbageCollectorControllerConfiguration(in, out, s)
}

func autoConvert_v1beta1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration(in *GenericControllerManagerConfiguration, out *config.GenericControllerManagerConfiguration, s conversion.Scope) error {
	if err := configv1beta1.Convert_v1beta1_MySQLConfiguration_To_config_MySQLConfiguration(&in.MySQL, &out.MySQL, s); err != nil {
		return err
	}
	if err := v1alpha1.Convert_v1alpha1_LeaderElectionConfiguration_To_config_LeaderElectionConfiguration(&in.LeaderElection, &out.LeaderElection, s); err != nil {
		return err
	}
	out.Namespace = in.Namespace
	out.BindAddress = in.BindAddress
	out.MetricsBindAddress = in.MetricsBindAddress
	out.HealthzBindAddress = in.HealthzBindAddress
	out.Parallelism = in.Parallelism
	out.SyncPeriod = in.SyncPeriod
	out.WatchFilterValue = in.WatchFilterValue
	return nil
}

// Convert_v1beta1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration is an autogenerated conversion function.
func Convert_v1beta1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration(in *GenericControllerManagerConfiguration, out *config.GenericControllerManagerConfiguration, s conversion.Scope) error {
	return autoConvert_v1beta1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration(in, out, s)
}

func autoConvert_config_GenericControllerManagerConfiguration_To_v1beta1_GenericControllerManagerConfiguration(in *config.GenericControllerManagerConfiguration, out *GenericControllerManagerConfiguration, s conversion.Scope) error {
	if err := configv1beta1.Convert_config_MySQLConfiguration_To_v1beta1_MySQLConfiguration(&in.MySQL, &out.MySQL, s); err != nil {
		return err
	}
	if err := v1alpha1.Convert_config_LeaderElectionConfiguration_To_v1alpha1_LeaderElectionConfiguration(&in.LeaderElection, &out.LeaderElection, s); err != nil {
		return err
	}
	out.Namespace = in.Namespace
	out.BindAddress = in.BindAddress
	out.MetricsBindAddress = in.MetricsBindAddress
	out.HealthzBindAddress = in.HealthzBindAddress
	out.Parallelism = in.Parallelism
	out.SyncPeriod = in.SyncPeriod
	out.WatchFilterValue = in.WatchFilterValue
	return nil
}

// Convert_config_GenericControllerManagerConfiguration_To_v1beta1_GenericControllerManagerConfiguration is an autogenerated conversion function.
func Convert_config_GenericControllerManagerConfiguration_To_v1beta1_GenericControllerManagerConfiguration(in *config.GenericControllerManagerConfiguration, out *GenericControllerManagerConfiguration, s conversion.Scope) error {
	return autoConvert_config_GenericControllerManagerConfiguration_To_v1beta1_GenericControllerManagerConfiguration(in, out, s)
}

func autoConvert_v1beta1_GroupResource_To_config_GroupResource(in *GroupResource, out *config.GroupResource, s conversion.Scope) error {
	out.Group = in.Group
	out.Resource = in.Resource
	return nil
}

// Convert_v1beta1_GroupResource_To_config_GroupResource is an autogenerated conversion function.
func Convert_v1beta1_GroupResource_To_config_GroupResource(in *GroupResource, out *config.GroupResource, s conversion.Scope) error {
	return autoConvert_v1beta1_GroupResource_To_config_GroupResource(in, out, s)
}

func autoConvert_config_GroupResource_To_v1beta1_GroupResource(in *config.GroupResource, out *GroupResource, s conversion.Scope) error {
	out.Group = in.Group
	out.Resource = in.Resource
	return nil
}

// Convert_config_GroupResource_To_v1beta1_GroupResource is an autogenerated conversion function.
func Convert_config_GroupResource_To_v1beta1_GroupResource(in *config.GroupResource, out *GroupResource, s conversion.Scope) error {
	return autoConvert_config_GroupResource_To_v1beta1_GroupResource(in, out, s)
}

func autoConvert_v1beta1_ZeroControllerManagerConfiguration_To_config_ZeroControllerManagerConfiguration(in *ZeroControllerManagerConfiguration, out *config.ZeroControllerManagerConfiguration, s conversion.Scope) error {
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	if err := Convert_v1beta1_GenericControllerManagerConfiguration_To_config_GenericControllerManagerConfiguration(&in.Generic, &out.Generic, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_GarbageCollectorControllerConfiguration_To_config_GarbageCollectorControllerConfiguration(&in.GarbageCollectorController, &out.GarbageCollectorController, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_ChainControllerConfiguration_To_config_ChainControllerConfiguration(&in.ChainController, &out.ChainController, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_ZeroControllerManagerConfiguration_To_config_ZeroControllerManagerConfiguration is an autogenerated conversion function.
func Convert_v1beta1_ZeroControllerManagerConfiguration_To_config_ZeroControllerManagerConfiguration(in *ZeroControllerManagerConfiguration, out *config.ZeroControllerManagerConfiguration, s conversion.Scope) error {
	return autoConvert_v1beta1_ZeroControllerManagerConfiguration_To_config_ZeroControllerManagerConfiguration(in, out, s)
}

func autoConvert_config_ZeroControllerManagerConfiguration_To_v1beta1_ZeroControllerManagerConfiguration(in *config.ZeroControllerManagerConfiguration, out *ZeroControllerManagerConfiguration, s conversion.Scope) error {
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	if err := Convert_config_GenericControllerManagerConfiguration_To_v1beta1_GenericControllerManagerConfiguration(&in.Generic, &out.Generic, s); err != nil {
		return err
	}
	if err := Convert_config_GarbageCollectorControllerConfiguration_To_v1beta1_GarbageCollectorControllerConfiguration(&in.GarbageCollectorController, &out.GarbageCollectorController, s); err != nil {
		return err
	}
	if err := Convert_config_ChainControllerConfiguration_To_v1beta1_ChainControllerConfiguration(&in.ChainController, &out.ChainController, s); err != nil {
		return err
	}
	return nil
}

// Convert_config_ZeroControllerManagerConfiguration_To_v1beta1_ZeroControllerManagerConfiguration is an autogenerated conversion function.
func Convert_config_ZeroControllerManagerConfiguration_To_v1beta1_ZeroControllerManagerConfiguration(in *config.ZeroControllerManagerConfiguration, out *ZeroControllerManagerConfiguration, s conversion.Scope) error {
	return autoConvert_config_ZeroControllerManagerConfiguration_To_v1beta1_ZeroControllerManagerConfiguration(in, out, s)
}