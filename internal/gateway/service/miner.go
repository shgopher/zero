// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package service

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/superproj/zero/internal/pkg/zerox"
	v1 "github.com/superproj/zero/pkg/api/gateway/v1"
	"github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

func (s *GatewayService) CreateMiner(ctx context.Context, m *v1beta1.Miner) (*emptypb.Empty, error) {
	if err := s.biz.Miners().Create(ctx, zerox.FromUserID(ctx), m); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *GatewayService) ListMiner(ctx context.Context, req *v1.ListMinerRequest) (*v1.ListMinerResponse, error) {
	ms, err := s.biz.Miners().List(ctx, zerox.FromUserID(ctx), req)
	if err != nil {
		return &v1.ListMinerResponse{}, err
	}

	return ms, nil
}

func (s *GatewayService) GetMiner(ctx context.Context, req *v1.GetMinerRequest) (*v1beta1.Miner, error) {
	m, err := s.biz.Miners().Get(ctx, zerox.FromUserID(ctx), req.Name)
	if err != nil {
		return &v1beta1.Miner{}, err
	}

	return m, nil
}

func (s *GatewayService) UpdateMiner(ctx context.Context, m *v1beta1.Miner) (*emptypb.Empty, error) {
	if err := s.biz.Miners().Update(ctx, zerox.FromUserID(ctx), m); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *GatewayService) DeleteMiner(ctx context.Context, req *v1.DeleteMinerRequest) (*emptypb.Empty, error) {
	if err := s.biz.Miners().Delete(ctx, zerox.FromUserID(ctx), req.Name); err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
