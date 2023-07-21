// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package mongo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/vinllen/mgo"
	"k8s.io/klog/v2"
)

type SinkConfig struct {
	CollectionName            string
	CollectionCapMaxSizeBytes int
	CollectionCapEnable       bool
}

// MongoSink represents an Mongo sink connector.
type MongoSink struct {
	ctx     context.Context
	conf    SinkConfig
	session *mgo.Session
	in      chan interface{}
}

// NewMongoSink returns a new MongoSink instance.
func NewMongoSink(ctx context.Context, session *mgo.Session, conf SinkConfig) (*MongoSink, error) {
	sink := &MongoSink{
		ctx:     ctx,
		conf:    conf,
		session: session,
		in:      make(chan interface{}),
	}

	go sink.init()
	return sink, nil
}

// init starts the main loop.
func (ms *MongoSink) init() {
	ms.capCollection()

	for msg := range ms.in {
		sess := ms.session.Copy()
		//nolint:staticcheck
		defer sess.Close()

		if err := sess.DB("").C(ms.conf.CollectionName).Insert(msg); err != nil {
			klog.ErrorS(err, "Problem inserting to mongo collection")
			if strings.Contains(strings.ToLower(err.Error()), "closed explicitly") {
				klog.V(1).InfoS("Detected connection failure!")
			}
		}
	}

	ms.session.Close()
}

func (ms *MongoSink) capCollection() (ok bool) {
	colName := ms.conf.CollectionName
	colCapMaxSizeBytes := ms.conf.CollectionCapMaxSizeBytes
	colCapEnable := ms.conf.CollectionCapEnable

	if !colCapEnable {
		return false
	}

	exists, err := ms.collectionExists(colName)
	if err != nil {
		klog.ErrorS(err, "Unable to determine if collection exists. Not capping collection", "collection", colName)
		return false
	}

	if exists {
		klog.V(1).InfoS("Collection already exists. Capping could result in data loss. Ignoring", "collection", colName)
		return false
	}

	if strconv.IntSize < 64 {
		klog.V(1).InfoS("Pump running < 64bit architecture. Not capping collection as max size would be 2gb")
		return false
	}

	if colCapMaxSizeBytes == 0 {
		defaultBytes := 5
		colCapMaxSizeBytes = defaultBytes * 1024 * 1024 * 1024
		klog.InfoS("No max collection size set for connection, set default value", "connection", colName, "size", colCapMaxSizeBytes)
	}

	sess := ms.session.Copy()
	defer sess.Close()

	err = ms.session.DB("").C(colName).Create(&mgo.CollectionInfo{Capped: true, MaxBytes: colCapMaxSizeBytes})
	if err != nil {
		klog.ErrorS(err, "Unable to create capped collection", "collection", colName)
		return false
	}

	klog.InfoS("Capped collection created", "collection", colName, "bytes", colCapMaxSizeBytes)

	return true
}

// collectionExists checks to see if a collection name exists in the db.
func (ms *MongoSink) collectionExists(name string) (bool, error) {
	sess := ms.session.Copy()
	defer sess.Close()

	colNames, err := sess.DB("").CollectionNames()
	if err != nil {
		klog.ErrorS(err, "Unable to get column names")
		return false, fmt.Errorf("failed to get collection name: %w", err)
	}

	for _, coll := range colNames {
		if coll == name {
			return true, nil
		}
	}

	return false, nil
}

// In returns an input channel for receiving data.
func (ks *MongoSink) In() chan<- interface{} {
	return ks.in
}
