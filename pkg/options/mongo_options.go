// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package options

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/spf13/pflag"
	"github.com/vinllen/mgo"
)

var _ IOptions = (*MongoOptions)(nil)

// MongoOptions defines options for mongo db.
type MongoOptions struct {
	// options needed when connnect to mongo db.
	TLSOptions *TLSOptions   `json:"tls" mapstructure:"tls"`
	URL        string        `json:"url" mapstructure:"url"`
	Timeout    time.Duration `json:"timeout" mapstructure:"timeout"`

	// mongo specific options
	CollectionName string `json:"collection-name"               mapstructure:"collection-name"`
}

// NewMongoOptions create a `zero` value instance.
func NewMongoOptions() *MongoOptions {
	return &MongoOptions{
		Timeout:    30 * time.Second,
		TLSOptions: NewTLSOptions(),
	}
}

// Validate verifies flags passed to MongoOptions.
func (o *MongoOptions) Validate() []error {
	errs := []error{}

	errs = append(errs, o.TLSOptions.Validate()...)

	return errs
}

// AddFlags adds flags related to redis storage for a specific APIServer to the specified FlagSet.
func (o *MongoOptions) AddFlags(fs *pflag.FlagSet, prefixs ...string) {
	o.TLSOptions.AddFlags(fs, "mongo")
	fs.DurationVar(&o.Timeout, "mongo.timeout", o.Timeout, "Timeout is the maximum amount of time a dial will wait for a connect to complete.")
	fs.StringVar(&o.URL, "mongo.url", o.URL, "The mongo server address.")
	fs.StringVar(&o.CollectionName, "mongo.collection-name", o.CollectionName, "The mongo collection name.")
}

func (o *MongoOptions) DialInfo() (dialInfo *mgo.DialInfo, err error) {
	if dialInfo, err = mgo.ParseURL(o.URL); err != nil {
		return dialInfo, fmt.Errorf("failed to parse mongo url: %w", err)
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		tlsConfig, err := o.TLSOptions.TLSConfig()
		if err != nil {
			return nil, err
		}

		return tls.Dial("tcp", addr.String(), tlsConfig)
	}
	dialInfo.Timeout = o.Timeout

	return dialInfo, err
}

func (o *MongoOptions) Session() (*mgo.Session, error) {
	dialInfo, err := o.DialInfo()
	if err != nil {
		return nil, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	return session, nil
}
