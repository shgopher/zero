// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package toyblc

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/superproj/zero/internal/pkg/core"
	"github.com/superproj/zero/internal/toyblc/blc"
	"github.com/superproj/zero/internal/toyblc/controller/v1/block"
	"github.com/superproj/zero/internal/toyblc/controller/v1/peer"
	"github.com/superproj/zero/internal/toyblc/ws"
	v1 "github.com/superproj/zero/pkg/api/toyblc/v1"
)

func installRouters(g *gin.Engine, bs *blc.BlockSet, ss *ws.Sockets) {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, v1.ErrorPageNotFound("route not found"), nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	// 注册 pprof 路由
	pprof.Register(g)

	bc := block.New(bs, ss)
	pc := peer.New(bs, ss)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 blocks 路由分组
		userv1 := v1.Group("/blocks")
		{
			userv1.POST("", bc.Create)
			userv1.GET("", bc.List)
		}

		// 创建 peers 路由分组
		postv1 := v1.Group("/peers")
		{
			postv1.POST("", pc.Create)
			postv1.GET("", pc.List)
		}
	}
}
