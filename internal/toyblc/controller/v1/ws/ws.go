// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package ws

import (
	"github.com/google/wire"
	"golang.org/x/net/websocket"

	"github.com/superproj/zero/internal/toyblc/blc"
	"github.com/superproj/zero/internal/toyblc/ws"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(New)

type WSController struct {
	bs *blc.BlockSet
	ss *ws.Sockets
}

func New(bs *blc.BlockSet, ss *ws.Sockets) *WSController {
	return &WSController{bs: bs, ss: ss}
}

func (wsc *WSController) WSHandler(w *websocket.Conn) {
	ws.WSHandler(wsc.bs, wsc.ss, w)
}
