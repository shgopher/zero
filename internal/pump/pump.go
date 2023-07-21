// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//go:generate wire .
package pump

import (
	"github.com/segmentio/kafka-go"
	"github.com/vinllen/mgo"

	contextutil "github.com/superproj/zero/internal/pkg/util/context"
	genericoptions "github.com/superproj/zero/pkg/options"
	kafkaconnector "github.com/superproj/zero/pkg/streams/connector/kafka"
	mongoconnector "github.com/superproj/zero/pkg/streams/connector/mongo"
	"github.com/superproj/zero/pkg/streams/flow"
)

// Config defines the config for the apiserver.
type Config struct {
	HealthOptions *genericoptions.HealthOptions
	KafkaOptions  *genericoptions.KafkaOptions
	MongoOptions  *genericoptions.MongoOptions
}

// Server contains state for a Kubernetes cluster master/api server.
type Server struct {
	config     kafka.ReaderConfig
	session    *mgo.Session
	collection string
}

type completedConfig struct {
	*Config
}

var addUTC = func(msg kafka.Message) kafka.Message {
	msg.Value = []byte("aaaaa")
	return msg
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

// New returns a new instance of Server from the given config.
// Certain config fields will be set to a default value if unset.
func (c completedConfig) New() (*Server, error) {
	session, err := c.MongoOptions.Session()
	if err != nil {
		return nil, err
	}

	server := &Server{
		config: kafka.ReaderConfig{
			Brokers:           c.KafkaOptions.Brokers,
			Topic:             c.KafkaOptions.Topic,
			GroupID:           c.KafkaOptions.ReaderOptions.GroupID,
			QueueCapacity:     c.KafkaOptions.ReaderOptions.QueueCapacity,
			MinBytes:          c.KafkaOptions.ReaderOptions.MinBytes,
			MaxBytes:          c.KafkaOptions.ReaderOptions.MaxBytes,
			MaxWait:           c.KafkaOptions.ReaderOptions.MaxWait,
			ReadBatchTimeout:  c.KafkaOptions.ReaderOptions.ReadBatchTimeout,
			HeartbeatInterval: c.KafkaOptions.ReaderOptions.HeartbeatInterval,
			CommitInterval:    c.KafkaOptions.ReaderOptions.CommitInterval,
			RebalanceTimeout:  c.KafkaOptions.ReaderOptions.RebalanceTimeout,
			StartOffset:       c.KafkaOptions.ReaderOptions.StartOffset,
			MaxAttempts:       c.KafkaOptions.ReaderOptions.MaxAttempts,
		},
		session:    session,
		collection: c.MongoOptions.CollectionName,
	}

	return server, nil
}

type preparedServer struct {
	*Server
}

func (s *Server) PrepareRun() preparedServer {
	return preparedServer{s}
}

func (s preparedServer) Run(stopCh <-chan struct{}) error {
	ctx, _ := contextutil.ContextForChannel(stopCh)

	source, err := kafkaconnector.NewKafkaSource(ctx, s.config)
	if err != nil {
		return err
	}

	filter := flow.NewMap(addUTC, 1)

	sink, err := mongoconnector.NewMongoSink(ctx, s.session, mongoconnector.SinkConfig{
		CollectionName:            s.collection,
		CollectionCapMaxSizeBytes: 5 * genericoptions.GiB,
		CollectionCapEnable:       true,
	})
	if err != nil {
		return err
	}
	source.Via(filter).To(sink)
	return nil
}
