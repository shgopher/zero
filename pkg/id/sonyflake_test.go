// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package id

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewSonyflake(t *testing.T) {
	sf1 := NewSonyflake(WithSonyflakeMachineId(1))
	sf2 := NewSonyflake(WithSonyflakeMachineId(2))
	sf3 := NewSonyflake(WithSonyflakeMachineId(3), WithSonyflakeStartTime(time.Date(1990, 10, 10, 0, 0, 0, 0, time.UTC)))
	i := 0
	for i < 100 {
		id1 := sf1.Id(context.Background())
		id2 := sf2.Id(context.Background())
		id3 := sf3.Id(context.Background())
		fmt.Println(id1, id2, id3)
		i++
	}
}
