// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package block

import (
	"github.com/gin-gonic/gin"

	"github.com/superproj/zero/internal/pkg/core"
	"github.com/superproj/zero/internal/toyblc/miner"
	v1 "github.com/superproj/zero/pkg/api/toyblc/v1"
)

func (b *BlockController) Create(c *gin.Context) {
	var r v1.CreateBlockRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	_ = miner.MinerBlock(b.bs, b.ss, r.Data)
	core.WriteResponse(c, nil, nil)
}
