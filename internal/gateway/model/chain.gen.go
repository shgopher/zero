// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameChainM = "chain"

// ChainM mapped from table <chain>
type ChainM struct {
	ID                     int64     `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true" json:"id"`
	Namespace              string    `gorm:"column:namespace;type:varchar(253);not null;uniqueIndex:uniq_namespace_name,priority:1" json:"namespace"`
	Name                   string    `gorm:"column:name;type:varchar(253);not null;uniqueIndex:uniq_namespace_name,priority:2" json:"name"`
	DisplayName            string    `gorm:"column:displayName;type:varchar(253);not null" json:"displayName"`
	MinerType              string    `gorm:"column:minerType;type:varchar(16);not null" json:"minerType"`
	Image                  string    `gorm:"column:image;type:varchar(253);not null" json:"image"`
	MinMineIntervalSeconds int32     `gorm:"column:minMineIntervalSeconds;type:int(8);not null" json:"minMineIntervalSeconds"`
	CreatedAt              time.Time `gorm:"column:createdAt;type:timestamp;not null;default:current_timestamp()" json:"createdAt"`
	UpdatedAt              time.Time `gorm:"column:updatedAt;type:timestamp;not null;default:current_timestamp()" json:"updatedAt"`
}

// TableName ChainM's table name
func (*ChainM) TableName() string {
	return TableNameChainM
}
