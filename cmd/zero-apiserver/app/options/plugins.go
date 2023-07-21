// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:staticcheck
package options

// This file exists to force the desired plugin implementations to be linked.
// This should probably be part of some configuration fed into the build for a
// given binary target.
import (
	// Admission policies.

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"

	"github.com/superproj/zero/internal/admission/plugin/admit"
	"github.com/superproj/zero/internal/admission/plugin/deny"
	"github.com/superproj/zero/internal/admission/plugin/namespace/autoprovision"
	"github.com/superproj/zero/internal/admission/plugin/namespace/exists"
	"github.com/superproj/zero/internal/admission/plugin/namespace/lifecycle"
)

// AllOrderedPlugins is the list of all the plugins in order.
var AllOrderedPlugins = []string{
	admit.PluginName,         // AlwaysAdmit
	autoprovision.PluginName, // NamespaceAutoProvision
	lifecycle.PluginName,     // NamespaceLifecycle
	exists.PluginName,        // NamespaceExists

	// new admission plugins should generally be inserted above here
	// webhook, resourcequota, and deny plugins must go at the end

	// mutatingwebhook.PluginName,   // MutatingAdmissionWebhook
	// validatingwebhook.PluginName, // ValidatingAdmissionWebhook
	// resourcequota.PluginName, // ResourceQuota
	deny.PluginName, // AlwaysDeny
}

// RegisterAllAdmissionPlugins registers all admission plugins.
// The order of registration is irrelevant, see AllOrderedPlugins for execution order.
func RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
	admit.Register(plugins) // DEPRECATED as no real meaning
	autoprovision.Register(plugins)
	lifecycle.Register(plugins)
	exists.Register(plugins)
	deny.Register(plugins) // DEPRECATED as no real meaning
}

// DefaultOffAdmissionPlugins get admission plugins off by default for zero-apiserver.
func DefaultOffAdmissionPlugins() sets.String {
	defaultOnPlugins := sets.NewString(
		autoprovision.PluginName, // NamespaceAutoProvision
		lifecycle.PluginName,     // NamespaceLifecycle
	)

	defaultOffPlugins := sets.NewString(AllOrderedPlugins...).Difference(defaultOnPlugins)

	return defaultOffPlugins
}
