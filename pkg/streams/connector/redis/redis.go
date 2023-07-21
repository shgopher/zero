// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/superproj/zero/pkg/streams"
	"github.com/superproj/zero/pkg/streams/flow"
)

// RedisSource represents a Redis Pub/Sub source connector.
type RedisSource struct {
	ctx     context.Context
	redisdb *redis.Client
	channel string
	out     chan interface{}
}

// NewRedisSource returns a new RedisSource instance.
func NewRedisSource(ctx context.Context, config *redis.Options, channel string) (*RedisSource, error) {
	redisdb := redis.NewClient(config)
	pubsub := redisdb.Subscribe(ctx, channel)

	// Wait for a confirmation that subscription is created before publishing anything
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	source := &RedisSource{
		ctx:     ctx,
		redisdb: redisdb,
		channel: channel,
		out:     make(chan interface{}),
	}

	go source.init(pubsub.Channel())
	return source, nil
}

// init starts the main loop.
func (rs *RedisSource) init(ch <-chan *redis.Message) {
loop:
	for {
		select {
		case <-rs.ctx.Done():
			break loop

		case msg := <-ch:
			rs.out <- msg
		}
	}

	log.Printf("Closing redis consumer")
	close(rs.out)
	rs.redisdb.Close()
}

// Via streams data through the given flow.
func (rs *RedisSource) Via(_flow streams.Flow) streams.Flow {
	flow.DoStream(rs, _flow)
	return _flow
}

// Out returns an output channel for sending data.
func (rs *RedisSource) Out() <-chan interface{} {
	return rs.out
}

// RedisSink represents a Redis Pub/Sub sink connector.
type RedisSink struct {
	redisdb *redis.Client
	channel string
	in      chan interface{}
}

// NewRedisSink returns a new RedisSink instance.
func NewRedisSink(config *redis.Options, channel string) *RedisSink {
	sink := &RedisSink{
		redisdb: redis.NewClient(config),
		channel: channel,
		in:      make(chan interface{}),
	}

	go sink.init()
	return sink
}

// init starts the main loop.
func (rs *RedisSink) init() {
	for msg := range rs.in {
		switch m := msg.(type) {
		case string:
			err := rs.redisdb.Publish(context.Background(), rs.channel, m).Err()
			if err != nil {
				log.Printf("redisdb.Publish failed with: %s", err)
			}

		default:
			log.Printf("Unsupported message type %v", m)
		}
	}

	log.Printf("Closing redis producer")
	rs.redisdb.Close()
}

// In returns an input channel for receiving data.
func (rs *RedisSink) In() chan<- interface{} {
	return rs.in
}
