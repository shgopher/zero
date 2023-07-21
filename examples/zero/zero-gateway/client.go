// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/superproj/zero/pkg/api/gateway/v1"
)

func main() {
	callHTTP()
	//callGRPC()
}

func callHTTP() {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("zero.gateway.superproj.com:51000"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewGatewayHTTPClient(conn)
	header := &http.Header{}
	header.Set("Authorization", "Bearer ccccc")

	reply, err := client.GetMinerSet(context.Background(), &v1.GetMinerSetRequest{Name: "test"}, transhttp.Header(header))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] GetMinerSet %s\n", reply.Spec.DisplayName)
	if errors.IsBadRequest(err) {
		log.Printf("[http] Login error is invalid argument: %v\n", err)
	}
}

func callGRPC() {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("zero.gateway.superproj.com:51010"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewGatewayClient(conn)
	reply, err := client.GetMinerSet(context.Background(), &v1.GetMinerSetRequest{Name: "test"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] GetMinerSet %+v\n", reply.Spec.DisplayName)
	if errors.IsBadRequest(err) {
		log.Printf("[grpc] Login error is invalid argument: %v\n", err)
	}
}
