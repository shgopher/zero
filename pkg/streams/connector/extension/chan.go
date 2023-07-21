// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package extension

import (
	"github.com/superproj/zero/pkg/streams"
	"github.com/superproj/zero/pkg/streams/flow"
)

// ChanSource represents an inbound connector that streams items from a channel.
type ChanSource struct {
	in chan interface{}
}

// NewChanSource returns a new ChanSource instance.
func NewChanSource(in chan interface{}) *ChanSource {
	return &ChanSource{in}
}

// Via streams data through the given flow.
func (cs *ChanSource) Via(_flow streams.Flow) streams.Flow {
	flow.DoStream(cs, _flow)
	return _flow
}

// Out returns an output channel for sending data.
func (cs *ChanSource) Out() <-chan interface{} {
	return cs.in
}

// ChanSink represents an outbound connector that streams items to a channel.
type ChanSink struct {
	Out chan interface{}
}

// NewChanSink returns a new ChanSink instance.
func NewChanSink(out chan interface{}) *ChanSink {
	return &ChanSink{out}
}

// In returns an input channel for receiving data.
func (ch *ChanSink) In() chan<- interface{} {
	return ch.Out
}
