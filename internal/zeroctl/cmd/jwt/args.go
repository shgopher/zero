// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package jwt

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ArgList defines a new pflag Value.
type ArgList map[string]string

// String return value of ArgList in string format.
func (l ArgList) String() string {
	data, _ := json.Marshal(l)

	return string(data)
}

// Set sets the value of ArgList.
func (l ArgList) Set(arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid argument '%v'. Must use format 'key=value'. %v", arg, parts)
	}
	l[parts[0]] = parts[1]

	return nil
}

// Type returns the type name of ArgList.
func (l ArgList) Type() string {
	return "map"
}
