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

const TableNameUserM = "user"

// UserM mapped from table <user>
type UserM struct {
	ID        int64     `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true" json:"id"`
	UserID    string    `gorm:"column:userID;type:varchar(253);not null;uniqueIndex:idx_userID,priority:1" json:"userID"`
	Username  string    `gorm:"column:username;type:varchar(253);not null;uniqueIndex:idx_username,priority:1" json:"username"`
	Status    int32     `gorm:"column:status;type:int(2);not null;default:1" json:"status"`
	Nickname  string    `gorm:"column:nickname;type:varchar(253);not null" json:"nickname"`
	Password  string    `gorm:"column:password;type:varchar(64);not null" json:"password"`
	Email     string    `gorm:"column:email;type:varchar(253);not null" json:"email"`
	Phone     string    `gorm:"column:phone;type:varchar(16);not null" json:"phone"`
	CreatedAt time.Time `gorm:"column:createdAt;type:timestamp;not null;default:current_timestamp()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:timestamp;not null;default:current_timestamp()" json:"updatedAt"`
}

// TableName UserM's table name
func (*UserM) TableName() string {
	return TableNameUserM
}
