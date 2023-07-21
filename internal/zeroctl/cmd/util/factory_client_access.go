// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// this file contains factories with no other dependencies

package util

import (
	"context"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/klog/v2"

	clioptions "github.com/superproj/zero/internal/zeroctl/util/options"
	usercenterv1 "github.com/superproj/zero/pkg/api/usercenter/v1"
)

type factoryImpl struct {
	opts *clioptions.Options
}

var _ Factory = (*factoryImpl)(nil)

func NewFactory(opts *clioptions.Options) Factory {
	if opts == nil {
		klog.Fatal("attempt to instantiate client_access_factory with nil clientGetter")
	}

	return &factoryImpl{opts: opts}
}

func (f *factoryImpl) GetOptions() *clioptions.Options {
	return f.opts
}

func (f *factoryImpl) UserCenterClient() usercenterv1.UserCenterClient {
	conn := creaeteConn(f.opts.UserCenterOptions)
	return usercenterv1.NewUserCenterClient(conn)
}

func (f *factoryImpl) MustWithToken(ctx context.Context) context.Context {
	ctx, err := f.WithToken(ctx)
	if err != nil {
		klog.Fatal(err.Error())
	}
	return ctx
}

func (f *factoryImpl) WithToken(ctx context.Context) (context.Context, error) {
	token, err := f.Login()
	if err != nil {
		return ctx, err
	}

	md := metadata.Pairs("Authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, nil
}

func (f *factoryImpl) Login() (token string, err error) {
	client := usercenterv1.NewUserCenterClient(creaeteConn(f.opts.UserCenterOptions))
	rp, err := client.Login(context.Background(), &usercenterv1.LoginRequest{
		Username: f.opts.UserOptions.Username,
		Password: f.opts.UserOptions.Password,
	})
	if err != nil {
		return "", err
	}

	klog.V(4).Infof("Get login token: %s", rp.Token)
	return rp.Token, nil
}

func creaeteConn(opts *clioptions.ServerOptions) grpc.ClientConnInterface {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint(opts.Addr),
		transgrpc.WithTimeout(opts.Timeout),
	)
	if err != nil {
		panic(err)
	}

	return conn
}
