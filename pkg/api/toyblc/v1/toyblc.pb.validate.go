// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: toyblc/v1/toyblc.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on CreateBlockRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateBlockRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateBlockRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateBlockRequestMultiError, or nil if none found.
func (m *CreateBlockRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateBlockRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Data

	if len(errors) > 0 {
		return CreateBlockRequestMultiError(errors)
	}

	return nil
}

// CreateBlockRequestMultiError is an error wrapping multiple validation errors
// returned by CreateBlockRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateBlockRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateBlockRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateBlockRequestMultiError) AllErrors() []error { return m }

// CreateBlockRequestValidationError is the validation error returned by
// CreateBlockRequest.Validate if the designated constraints aren't met.
type CreateBlockRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateBlockRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateBlockRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateBlockRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateBlockRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateBlockRequestValidationError) ErrorName() string {
	return "CreateBlockRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateBlockRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateBlockRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateBlockRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateBlockRequestValidationError{}

// Validate checks the field values on CreatePeerRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CreatePeerRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePeerRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePeerRequestMultiError, or nil if none found.
func (m *CreatePeerRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePeerRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Peer

	if len(errors) > 0 {
		return CreatePeerRequestMultiError(errors)
	}

	return nil
}

// CreatePeerRequestMultiError is an error wrapping multiple validation errors
// returned by CreatePeerRequest.ValidateAll() if the designated constraints
// aren't met.
type CreatePeerRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePeerRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePeerRequestMultiError) AllErrors() []error { return m }

// CreatePeerRequestValidationError is the validation error returned by
// CreatePeerRequest.Validate if the designated constraints aren't met.
type CreatePeerRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePeerRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePeerRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePeerRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePeerRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePeerRequestValidationError) ErrorName() string {
	return "CreatePeerRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePeerRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePeerRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePeerRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePeerRequestValidationError{}
