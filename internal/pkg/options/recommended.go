// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package options

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	admissionmetrics "k8s.io/apiserver/pkg/admission/metrics"
	"k8s.io/apiserver/pkg/authentication/authenticatorfactory"
	"k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/util/feature"
	utilflowcontrol "k8s.io/apiserver/pkg/util/flowcontrol"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	openapicommon "k8s.io/kube-openapi/pkg/common"
)

var configScheme = runtime.NewScheme()

// RecommendedOptions contains the recommended options for running an API server.
// If you add something to this list, it should be in a logical grouping.
// Each of them can be nil to leave the feature unconfigured on ApplyTo.
type RecommendedOptions struct {
	*genericoptions.RecommendedOptions
}

func NewRecommendedOptions(prefix string, codec runtime.Codec) *RecommendedOptions {
	return &RecommendedOptions{genericoptions.NewRecommendedOptions(prefix, codec)}
}

// ApplyTo adds RecommendedOptions to the server configuration.
// pluginInitializers can be empty, it is only need for additional initializers.
func (o *RecommendedOptions) ApplyTo(config *server.RecommendedConfig) error {
	if err := o.Etcd.ApplyTo(&config.Config); err != nil {
		return err
	}
	if err := o.EgressSelector.ApplyTo(&config.Config); err != nil {
		return err
	}
	if feature.DefaultFeatureGate.Enabled(features.APIServerTracing) {
		if err := o.Traces.ApplyTo(config.Config.EgressSelector, &config.Config); err != nil {
			return err
		}
	}
	if err := o.SecureServing.ApplyTo(&config.Config.SecureServing, &config.Config.LoopbackClientConfig); err != nil {
		return err
	}
	if err := authenticationApplyTo(o.Authentication, &config.Config.Authentication, config.SecureServing, config.OpenAPIConfig); err != nil {
		return err
	}
	if err := o.Authorization.ApplyTo(&config.Config.Authorization); err != nil {
		return err
	}
	if err := o.Audit.ApplyTo(&config.Config); err != nil {
		return err
	}
	if err := o.Features.ApplyTo(&config.Config); err != nil {
		return err
	}
	if err := o.CoreAPI.ApplyTo(config); err != nil {
		return err
	}
	if initializers, err := o.ExtraAdmissionInitializers(config); err != nil {
		return err
		// } else if err := admissionOptionsApplyTo(o.Admission, &config.Config, o.FeatureGate, initializers...); err != nil {
	} else if err := admissionOptionsApplyTo(o.Admission, &config.Config, initializers...); err != nil {
		return err
	}
	//nolint: nestif
	if feature.DefaultFeatureGate.Enabled(features.APIPriorityAndFairness) {
		if config.ClientConfig != nil {
			if config.MaxRequestsInFlight+config.MaxMutatingRequestsInFlight <= 0 {
				return fmt.Errorf(
					"invalid configuration: MaxRequestsInFlight=%d and MaxMutatingRequestsInFlight=%d; they must add up to something positive",
					config.MaxRequestsInFlight,
					config.MaxMutatingRequestsInFlight,
				)
			}
			config.FlowControl = utilflowcontrol.New(
				config.SharedInformerFactory,
				kubernetes.NewForConfigOrDie(config.ClientConfig).FlowcontrolV1beta3(),
				config.MaxRequestsInFlight+config.MaxMutatingRequestsInFlight,
				config.RequestTimeout/4,
			)
		} else {
			klog.Warningf("Neither kubeconfig is provided nor service-account is mounted, so APIPriorityAndFairness will be disabled")
		}
	}
	return nil
}

// admissionOptionsApplyTo adds the admission chain to the server configuration.
// In case admission plugin names were not provided by a cluster-admin they will be prepared from the
// recommended/default values.
// In addition the method lazily initializes a generic plugin that is appended to the list of pluginInitializers
// note this method uses:
//
//	genericconfig.Authorizer
func admissionOptionsApplyTo(
	a *genericoptions.AdmissionOptions,
	c *server.Config,
	// features featuregate.FeatureGate,
	pluginInitializers ...admission.PluginInitializer,
) error {
	if a == nil {
		return nil
	}

	pluginNames := enabledPluginNames(a)

	pluginsConfigProvider, err := admission.ReadAdmissionConfiguration(pluginNames, a.ConfigFile, configScheme)
	if err != nil {
		return fmt.Errorf("failed to read plugin config: %w", err)
	}

	initializersChain := admission.PluginInitializers{}
	initializersChain = append(initializersChain, pluginInitializers...)

	admissionChain, err := a.Plugins.NewFromPlugins(pluginNames, pluginsConfigProvider, initializersChain, a.Decorators)
	if err != nil {
		return err
	}

	c.AdmissionControl = admissionmetrics.WithStepMetrics(admissionChain)
	return nil
}

// enabledPluginNames makes use of RecommendedPluginOrder, DefaultOffPlugins,
// EnablePlugins, DisablePlugins fields
// to prepare a list of ordered plugin names that are enabled.
func enabledPluginNames(a *genericoptions.AdmissionOptions) []string {
	allOffPlugins := append(a.DefaultOffPlugins.List(), a.DisablePlugins...)
	disabledPlugins := sets.NewString(allOffPlugins...)
	enabledPlugins := sets.NewString(a.EnablePlugins...)
	disabledPlugins = disabledPlugins.Difference(enabledPlugins)
	orderedPlugins := []string{}
	for _, plugin := range a.RecommendedPluginOrder {
		if !disabledPlugins.Has(plugin) {
			orderedPlugins = append(orderedPlugins, plugin)
		}
	}

	return orderedPlugins
}

func authenticationApplyTo(s *genericoptions.DelegatingAuthenticationOptions, authenticationInfo *server.AuthenticationInfo,
	servingInfo *server.SecureServingInfo, openAPIConfig *openapicommon.Config,
) error {
	if s == nil {
		authenticationInfo.Authenticator = nil
		return nil
	}

	cfg := authenticatorfactory.DelegatingAuthenticatorConfig{
		Anonymous:                true,
		CacheTTL:                 s.CacheTTL,
		WebhookRetryBackoff:      s.WebhookRetryBackoff,
		TokenAccessReviewTimeout: s.TokenRequestTimeout,
	}

	var err error

	// get the clientCA information
	clientCASpecified := s.ClientCert != genericoptions.ClientCertAuthenticationOptions{}
	var clientCAProvider dynamiccertificates.CAContentProvider
	if clientCASpecified {
		clientCAProvider, err = s.ClientCert.GetClientCAContentProvider()
		if err != nil {
			return fmt.Errorf("unable to load client CA provider: %w", err)
		}
		cfg.ClientCertificateCAContentProvider = clientCAProvider
		if err = authenticationInfo.ApplyClientCert(cfg.ClientCertificateCAContentProvider, servingInfo); err != nil {
			return fmt.Errorf("unable to assign client CA provider: %w", err)
		}
	}

	requestHeaderCAFileSpecified := len(s.RequestHeader.ClientCAFile) > 0
	var requestHeaderConfig *authenticatorfactory.RequestHeaderConfig
	if requestHeaderCAFileSpecified {
		requestHeaderConfig, err = s.RequestHeader.ToAuthenticationRequestHeaderConfig()
		if err != nil {
			return fmt.Errorf("unable to create request header authentication config: %w", err)
		}
	}

	if requestHeaderConfig != nil {
		cfg.RequestHeaderConfig = requestHeaderConfig
		if err = authenticationInfo.ApplyClientCert(cfg.RequestHeaderConfig.CAContentProvider, servingInfo); err != nil {
			return fmt.Errorf("unable to load request-header-client-ca-file: %w", err)
		}
	}

	// create authenticator
	authenticator, securityDefinitions, err := cfg.New()
	if err != nil {
		return err
	}
	authenticationInfo.Authenticator = authenticator
	if openAPIConfig != nil {
		openAPIConfig.SecurityDefinitions = securityDefinitions
	}

	return nil
}
