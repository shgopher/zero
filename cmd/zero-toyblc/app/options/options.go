// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package options contains flags and options for initializing an apiserver
package options

import (
	"fmt"
	"time"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/superproj/zero/internal/toyblc"
	"github.com/superproj/zero/internal/toyblc/defaults"
	"github.com/superproj/zero/pkg/app"
	"github.com/superproj/zero/pkg/log"
	genericoptions "github.com/superproj/zero/pkg/options"
)

const (
	// UserAgent is the userAgent name when starting zero-gateway server.
	UserAgent = "zero-toyblc"
)

var _ app.CliOptions = (*Options)(nil)

// Options contains state for master/api server.
type Options struct {
	Miner            bool                        `json:"miner" mapstructure:"miner"`
	MinMineInterval  time.Duration               `json:"min-mine-interval" mapstructure:"min-mine-interval"`
	MiningDifficulty int                         `json:"mining-difficulty" mapstructure:"mining-difficulty"`
	Address          string                      `json:"account" mapstructure:"account"`
	P2PAddr          string                      `json:"p2p-addr" mapstructure:"p2p-addr"`
	Peers            []string                    `json:"peers" mapstructure:"peers"`
	HTTPOptions      *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
	Log              *log.Options                `json:"log" mapstructure:"log"`
}

// NewOptions returns initialized Options.
func NewOptions() *Options {
	o := &Options{
		MinMineInterval:  2 * time.Hour,
		MiningDifficulty: 1,
		Address:          defaults.GenesisAddress,
		P2PAddr:          "0.0.0.0:6001",
		Peers:            []string{"ws://localhost:6001"},
		HTTPOptions:      genericoptions.NewHTTPOptions(),
		Log:              log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.HTTPOptions.AddFlags(fss.FlagSet("http"))
	o.Log.AddFlags(fss.FlagSet("log"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.BoolVar(&o.Miner, "miner", o.Miner, "Turn on mining mode.")
	fs.DurationVar(&o.MinMineInterval, "min-mine-interval", o.MinMineInterval, "Specify the minimum mining interval.")
	fs.IntVar(&o.MiningDifficulty, "mining-difficulty", o.MiningDifficulty, "Specify the mining difficulty.")
	fs.StringVar(&o.Address, "address", o.Address, "Wallet account to receive the block rewards.")
	fs.StringVar(&o.P2PAddr, "p2p-addr", o.P2PAddr, "The p2p server address.")
	fs.StringSliceVar(&o.Peers, "peers", o.Peers, "The initial peers.")

	return fss
}

// Complete completes all the required options.
func (o *Options) Complete() error {
	return nil
}

// Validate validates all the required options.
func (o *Options) Validate() error {
	errs := []error{}

	if o.MiningDifficulty < 0 {
		errs = append(errs, fmt.Errorf("`--mining-difficulty` must be non-negative"))
	}

	if err := genericoptions.ValidateAddress(o.P2PAddr); err != nil {
		errs = append(errs, err)
	}

	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up zero-nightwatch config with options.
func (o *Options) ApplyTo(c *toyblc.Config) error {
	c.Miner = o.Miner
	c.MinMineInterval = o.MinMineInterval
	c.Address = o.Address
	c.HTTPOptions = o.HTTPOptions
	c.P2PAddr = o.P2PAddr
	c.Peers = o.Peers

	return nil
}

// Config return a zero-nightwatch config object.
func (o *Options) Config() (*toyblc.Config, error) {
	c := &toyblc.Config{}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
