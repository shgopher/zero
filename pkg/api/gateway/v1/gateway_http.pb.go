// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.6.1
// - protoc             v3.21.6
// source: gateway/v1/gateway.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	v1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationGatewayCreateMiner = "/gateway.v1.Gateway/CreateMiner"
const OperationGatewayCreateMinerSet = "/gateway.v1.Gateway/CreateMinerSet"
const OperationGatewayDeleteMiner = "/gateway.v1.Gateway/DeleteMiner"
const OperationGatewayDeleteMinerSet = "/gateway.v1.Gateway/DeleteMinerSet"
const OperationGatewayGetIdempotentToken = "/gateway.v1.Gateway/GetIdempotentToken"
const OperationGatewayGetMiner = "/gateway.v1.Gateway/GetMiner"
const OperationGatewayGetMinerSet = "/gateway.v1.Gateway/GetMinerSet"
const OperationGatewayListMiner = "/gateway.v1.Gateway/ListMiner"
const OperationGatewayListMinerSet = "/gateway.v1.Gateway/ListMinerSet"
const OperationGatewayScaleMinerSet = "/gateway.v1.Gateway/ScaleMinerSet"
const OperationGatewayUpdateMiner = "/gateway.v1.Gateway/UpdateMiner"
const OperationGatewayUpdateMinerSet = "/gateway.v1.Gateway/UpdateMinerSet"

type GatewayHTTPServer interface {
	// CreateMiner CreateMiner
	CreateMiner(context.Context, *v1beta1.Miner) (*emptypb.Empty, error)
	// CreateMinerSet CreateMinerSet
	CreateMinerSet(context.Context, *v1beta1.MinerSet) (*emptypb.Empty, error)
	// DeleteMiner DeleteMiner
	DeleteMiner(context.Context, *DeleteMinerRequest) (*emptypb.Empty, error)
	// DeleteMinerSet DeleteMinerSet
	DeleteMinerSet(context.Context, *DeleteMinerSetRequest) (*emptypb.Empty, error)
	// GetIdempotentToken GetIdempotentToken
	GetIdempotentToken(context.Context, *emptypb.Empty) (*IdempotentResponse, error)
	// GetMiner GetMiner
	GetMiner(context.Context, *GetMinerRequest) (*v1beta1.Miner, error)
	// GetMinerSet GetMinerSet
	GetMinerSet(context.Context, *GetMinerSetRequest) (*v1beta1.MinerSet, error)
	// ListMiner ListMiner
	ListMiner(context.Context, *ListMinerRequest) (*ListMinerResponse, error)
	// ListMinerSet ListMinerSet
	ListMinerSet(context.Context, *ListMinerSetRequest) (*ListMinerSetResponse, error)
	// ScaleMinerSet ScaleMinerSet
	ScaleMinerSet(context.Context, *ScaleMinerSetRequest) (*emptypb.Empty, error)
	// UpdateMiner UpdateMiner
	UpdateMiner(context.Context, *v1beta1.Miner) (*emptypb.Empty, error)
	// UpdateMinerSet UpdateMinerSet
	UpdateMinerSet(context.Context, *v1beta1.MinerSet) (*emptypb.Empty, error)
}

func RegisterGatewayHTTPServer(s *http.Server, srv GatewayHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/idempotents", _Gateway_GetIdempotentToken0_HTTP_Handler(srv))
	r.POST("/v1/minersets", _Gateway_CreateMinerSet0_HTTP_Handler(srv))
	r.GET("/v1/minersets", _Gateway_ListMinerSet0_HTTP_Handler(srv))
	r.GET("/v1/minersets/{name}", _Gateway_GetMinerSet0_HTTP_Handler(srv))
	r.PUT("/v1/minersets", _Gateway_UpdateMinerSet0_HTTP_Handler(srv))
	r.DELETE("/v1/minersets/{name}", _Gateway_DeleteMinerSet0_HTTP_Handler(srv))
	r.PUT("/v1/minersets/{name}/scale", _Gateway_ScaleMinerSet0_HTTP_Handler(srv))
	r.POST("/v1/miners", _Gateway_CreateMiner0_HTTP_Handler(srv))
	r.GET("/v1/miners", _Gateway_ListMiner0_HTTP_Handler(srv))
	r.GET("/v1/miners/{name}", _Gateway_GetMiner0_HTTP_Handler(srv))
	r.PUT("/v1/miners", _Gateway_UpdateMiner0_HTTP_Handler(srv))
	r.DELETE("/v1/miners/{name}", _Gateway_DeleteMiner0_HTTP_Handler(srv))
}

func _Gateway_GetIdempotentToken0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayGetIdempotentToken)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetIdempotentToken(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*IdempotentResponse)
		return ctx.Result(200, reply)
	}
}

func _Gateway_CreateMinerSet0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in v1beta1.MinerSet
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayCreateMinerSet)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateMinerSet(ctx, req.(*v1beta1.MinerSet))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Gateway_ListMinerSet0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListMinerSetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayListMinerSet)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListMinerSet(ctx, req.(*ListMinerSetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListMinerSetResponse)
		return ctx.Result(200, reply)
	}
}

func _Gateway_GetMinerSet0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetMinerSetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayGetMinerSet)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetMinerSet(ctx, req.(*GetMinerSetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*v1beta1.MinerSet)
		return ctx.Result(200, reply)
	}
}

func _Gateway_UpdateMinerSet0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in v1beta1.MinerSet
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayUpdateMinerSet)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateMinerSet(ctx, req.(*v1beta1.MinerSet))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Gateway_DeleteMinerSet0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteMinerSetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayDeleteMinerSet)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteMinerSet(ctx, req.(*DeleteMinerSetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Gateway_ScaleMinerSet0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ScaleMinerSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayScaleMinerSet)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ScaleMinerSet(ctx, req.(*ScaleMinerSetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Gateway_CreateMiner0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in v1beta1.Miner
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayCreateMiner)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateMiner(ctx, req.(*v1beta1.Miner))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Gateway_ListMiner0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListMinerRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayListMiner)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListMiner(ctx, req.(*ListMinerRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListMinerResponse)
		return ctx.Result(200, reply)
	}
}

func _Gateway_GetMiner0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetMinerRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayGetMiner)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetMiner(ctx, req.(*GetMinerRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*v1beta1.Miner)
		return ctx.Result(200, reply)
	}
}

func _Gateway_UpdateMiner0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in v1beta1.Miner
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayUpdateMiner)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateMiner(ctx, req.(*v1beta1.Miner))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Gateway_DeleteMiner0_HTTP_Handler(srv GatewayHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteMinerRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGatewayDeleteMiner)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteMiner(ctx, req.(*DeleteMinerRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type GatewayHTTPClient interface {
	CreateMiner(ctx context.Context, req *v1beta1.Miner, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	CreateMinerSet(ctx context.Context, req *v1beta1.MinerSet, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteMiner(ctx context.Context, req *DeleteMinerRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteMinerSet(ctx context.Context, req *DeleteMinerSetRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetIdempotentToken(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *IdempotentResponse, err error)
	GetMiner(ctx context.Context, req *GetMinerRequest, opts ...http.CallOption) (rsp *v1beta1.Miner, err error)
	GetMinerSet(ctx context.Context, req *GetMinerSetRequest, opts ...http.CallOption) (rsp *v1beta1.MinerSet, err error)
	ListMiner(ctx context.Context, req *ListMinerRequest, opts ...http.CallOption) (rsp *ListMinerResponse, err error)
	ListMinerSet(ctx context.Context, req *ListMinerSetRequest, opts ...http.CallOption) (rsp *ListMinerSetResponse, err error)
	ScaleMinerSet(ctx context.Context, req *ScaleMinerSetRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateMiner(ctx context.Context, req *v1beta1.Miner, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateMinerSet(ctx context.Context, req *v1beta1.MinerSet, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type GatewayHTTPClientImpl struct {
	cc *http.Client
}

func NewGatewayHTTPClient(client *http.Client) GatewayHTTPClient {
	return &GatewayHTTPClientImpl{client}
}

func (c *GatewayHTTPClientImpl) CreateMiner(ctx context.Context, in *v1beta1.Miner, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/miners"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGatewayCreateMiner))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) CreateMinerSet(ctx context.Context, in *v1beta1.MinerSet, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/minersets"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGatewayCreateMinerSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) DeleteMiner(ctx context.Context, in *DeleteMinerRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/miners/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGatewayDeleteMiner))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) DeleteMinerSet(ctx context.Context, in *DeleteMinerSetRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/minersets/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGatewayDeleteMinerSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) GetIdempotentToken(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*IdempotentResponse, error) {
	var out IdempotentResponse
	pattern := "/v1/idempotents"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGatewayGetIdempotentToken))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) GetMiner(ctx context.Context, in *GetMinerRequest, opts ...http.CallOption) (*v1beta1.Miner, error) {
	var out v1beta1.Miner
	pattern := "/v1/miners/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGatewayGetMiner))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) GetMinerSet(ctx context.Context, in *GetMinerSetRequest, opts ...http.CallOption) (*v1beta1.MinerSet, error) {
	var out v1beta1.MinerSet
	pattern := "/v1/minersets/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGatewayGetMinerSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) ListMiner(ctx context.Context, in *ListMinerRequest, opts ...http.CallOption) (*ListMinerResponse, error) {
	var out ListMinerResponse
	pattern := "/v1/miners"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGatewayListMiner))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) ListMinerSet(ctx context.Context, in *ListMinerSetRequest, opts ...http.CallOption) (*ListMinerSetResponse, error) {
	var out ListMinerSetResponse
	pattern := "/v1/minersets"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGatewayListMinerSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) ScaleMinerSet(ctx context.Context, in *ScaleMinerSetRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/minersets/{name}/scale"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGatewayScaleMinerSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) UpdateMiner(ctx context.Context, in *v1beta1.Miner, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/miners"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGatewayUpdateMiner))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *GatewayHTTPClientImpl) UpdateMinerSet(ctx context.Context, in *v1beta1.MinerSet, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/minersets"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGatewayUpdateMinerSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
