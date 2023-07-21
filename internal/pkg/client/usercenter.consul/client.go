// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package usercenter

import (
	"context"
)

func (c *clientImpl) GetSecret(ctx context.Context, req *GetSecretRequest) (*GetSecretResponse, error) {
	/*
		ret := &GetSecretResponse{}

		_, err := impl.r.SetPathParams(map[string]string{
			"name": req.Name,
			if err != nil {
				return nil, err
			}
		}).SetResult(ret).Get("/v1/secrets/{name}")

		return ret, err
	*/
	return nil, nil
}
