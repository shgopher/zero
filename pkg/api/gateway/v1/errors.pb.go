// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.6
// source: gateway/v1/errors.proto

package v1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ErrorReason int32

const (
	// user
	ErrorReason_UserLoginFailed   ErrorReason = 0
	ErrorReason_UserAlreadyExists ErrorReason = 1
	ErrorReason_UserNotFound      ErrorReason = 2
	ErrorReason_UserCreateFailed  ErrorReason = 3
)

// Enum value maps for ErrorReason.
var (
	ErrorReason_name = map[int32]string{
		0: "UserLoginFailed",
		1: "UserAlreadyExists",
		2: "UserNotFound",
		3: "UserCreateFailed",
	}
	ErrorReason_value = map[string]int32{
		"UserLoginFailed":   0,
		"UserAlreadyExists": 1,
		"UserNotFound":      2,
		"UserCreateFailed":  3,
	}
)

func (x ErrorReason) Enum() *ErrorReason {
	p := new(ErrorReason)
	*p = x
	return p
}

func (x ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_gateway_v1_errors_proto_enumTypes[0].Descriptor()
}

func (ErrorReason) Type() protoreflect.EnumType {
	return &file_gateway_v1_errors_proto_enumTypes[0]
}

func (x ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorReason.Descriptor instead.
func (ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_gateway_v1_errors_proto_rawDescGZIP(), []int{0}
}

var File_gateway_v1_errors_proto protoreflect.FileDescriptor

var file_gateway_v1_errors_proto_rawDesc = []byte{
	0x0a, 0x17, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x7f, 0x0a, 0x0b, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x0f, 0x55, 0x73, 0x65,
	0x72, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x00, 0x1a, 0x04,
	0xa8, 0x45, 0x91, 0x03, 0x12, 0x1b, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x41, 0x6c, 0x72, 0x65,
	0x61, 0x64, 0x79, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x10, 0x01, 0x1a, 0x04, 0xa8, 0x45, 0x99,
	0x03, 0x12, 0x16, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e,
	0x64, 0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03, 0x12, 0x1a, 0x0a, 0x10, 0x55, 0x73, 0x65,
	0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x03, 0x1a,
	0x04, 0xa8, 0x45, 0x9d, 0x04, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x33, 0x50, 0x01, 0x5a,
	0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x72, 0x6d,
	0x6f, 0x74, 0x65, 0x64, 0x75, 0x2f, 0x7a, 0x65, 0x72, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gateway_v1_errors_proto_rawDescOnce sync.Once
	file_gateway_v1_errors_proto_rawDescData = file_gateway_v1_errors_proto_rawDesc
)

func file_gateway_v1_errors_proto_rawDescGZIP() []byte {
	file_gateway_v1_errors_proto_rawDescOnce.Do(func() {
		file_gateway_v1_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_gateway_v1_errors_proto_rawDescData)
	})
	return file_gateway_v1_errors_proto_rawDescData
}

var file_gateway_v1_errors_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gateway_v1_errors_proto_goTypes = []interface{}{
	(ErrorReason)(0), // 0: gateway.v1.ErrorReason
}
var file_gateway_v1_errors_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gateway_v1_errors_proto_init() }
func file_gateway_v1_errors_proto_init() {
	if File_gateway_v1_errors_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gateway_v1_errors_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gateway_v1_errors_proto_goTypes,
		DependencyIndexes: file_gateway_v1_errors_proto_depIdxs,
		EnumInfos:         file_gateway_v1_errors_proto_enumTypes,
	}.Build()
	File_gateway_v1_errors_proto = out.File
	file_gateway_v1_errors_proto_rawDesc = nil
	file_gateway_v1_errors_proto_goTypes = nil
	file_gateway_v1_errors_proto_depIdxs = nil
}
