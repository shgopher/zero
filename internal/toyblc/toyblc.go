// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package toyblc

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"

	"github.com/superproj/zero/internal/toyblc/blc"
	wscontroller "github.com/superproj/zero/internal/toyblc/controller/v1/ws"
	mw "github.com/superproj/zero/internal/toyblc/middleware"
	"github.com/superproj/zero/internal/toyblc/miner"
	"github.com/superproj/zero/internal/toyblc/ws"
	"github.com/superproj/zero/pkg/log"
	genericoptions "github.com/superproj/zero/pkg/options"
)

// Config represents the configuration of the service.
type Config struct {
	Miner           bool
	MinMineInterval time.Duration
	Address         string
	HTTPOptions     *genericoptions.HTTPOptions
	P2PAddr         string
	Peers           []string
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

type completedConfig struct {
	*Config
}

// New returns a new instance of Server from the given config.
func (c completedConfig) New() (*Server, error) {
	bs, ss := blc.NewBlockSet(c.Address), ws.NewSockets()

	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure}

	// 创建 Gin 引擎
	g := gin.New()

	// 并初始化路由
	installRouters(g, bs, ss)

	// 添加中间件
	g.Use(mws...)

	// http.Handle
	// 创建 HTTP Server 实例
	httpsrv := &http.Server{Addr: c.HTTPOptions.Addr, Handler: g}

	p2p := gin.New()
	wsc := wscontroller.New(bs, ss)
	p2p.Use(gin.WrapH(websocket.Handler(wsc.WSHandler)))

	p2psrv := &http.Server{Addr: c.P2PAddr, Handler: p2p}
	return &Server{
		srv:             httpsrv,
		p2psrv:          p2psrv,
		bs:              bs,
		ss:              ss,
		miner:           c.Miner,
		minMineInterval: c.MinMineInterval,
		peers:           c.Peers,
	}, nil
}

// Server represents the server.
type Server struct {
	srv             *http.Server
	p2psrv          *http.Server
	bs              *blc.BlockSet
	ss              *ws.Sockets
	miner           bool
	minMineInterval time.Duration
	peers           []string
}

func (s *Server) Run(stopCh <-chan struct{}) error {
	if s.miner {
		miner.NewMiner(s.bs, s.ss, s.minMineInterval).Start()
	}

	// 运行 HTTP 服务器。在 goroutine 中启动服务器，它不会阻止下面的正常关闭处理流程
	// 打印一条日志，用来提示 HTTP 服务已经起来，方便排障
	// log.Infow("Start to listening the incoming requests on http address", "addr", "111")
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	// log.Infow("Start to listening the incoming requests on p2p address", "addr", "11")
	go func() {
		if err := s.p2psrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	ws.ConnectToPeers(s.bs, s.ss, s.peers)

	<-stopCh
	log.Infow("Shutting down server ...")

	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw(err, "Insecure Server forced to shutdown")
		return err
	}

	if err := s.p2psrv.Shutdown(ctx); err != nil {
		log.Errorw(err, "P2P Server forced to shutdown")
		return err
	}

	log.Infow("Server exiting")
	return nil
}
