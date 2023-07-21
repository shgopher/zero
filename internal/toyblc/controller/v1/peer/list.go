// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package peer

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/superproj/zero/internal/pkg/core"
)

func (b *PeerController) List(c *gin.Context) {
	var slice []string

	for _, socket := range b.ss.List() {
		if socket.IsClientConn() {
			slice = append(slice, strings.Replace(socket.LocalAddr().String(), "ws://", "", 1))
		} else {
			slice = append(slice, socket.Request().RemoteAddr)
		}
	}

	core.WriteResponse(c, nil, slice)
}
