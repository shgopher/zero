// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package options

import (
	"context"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var _ IOptions = (*JaegerOptions)(nil)

// JaegerOptions defines options for consul client.
type JaegerOptions struct {
	// Server is the url of the Jaeger server
	Server      string `json:"server,omitempty" mapstructure:"server"`
	ServiceName string `json:"service-name,omitempty" mapstructure:"service-name"`
	Env         string `json:"env,omitempty" mapstructure:"env"`
}

// NewJaegerOptions create a `zero` value instance.
func NewJaegerOptions() *JaegerOptions {
	return &JaegerOptions{
		Server: "http://127.0.0.1:14268/api/traces",
		Env:    "dev",
	}
}

// Validate verifies flags passed to JaegerOptions.
func (o *JaegerOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *JaegerOptions) AddFlags(fs *pflag.FlagSet, prefixs ...string) {
	fs.StringVar(&o.Server, "jaeger.server", o.Server, ""+
		"Server is the url of the Jaeger server.")
	fs.StringVar(&o.ServiceName, "jaeger.service-name", o.ServiceName, ""+
		"Specify the service name for jaeger resource.")
	fs.StringVar(&o.Env, "jaeger.env", o.Env, "Specify the deployment environment(dev/test/staging/prod).")
}

func (o *JaegerOptions) SetTracerProvider() error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Server)))
	if err != nil {
		return err
	}

	res, _ := resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceNameKey.String(o.ServiceName),
		attribute.String("env", o.Env),
		attribute.String("exporter", "jaeger"),
	))

	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return nil
}
