// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package conditions

import (
	"errors"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"

	appsv1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

// HaveSameStateOf matches a condition to have the same state of another.
func HaveSameStateOf(expected *appsv1beta1.Condition) types.GomegaMatcher {
	return &conditionMatcher{
		Expected: expected,
	}
}

type conditionMatcher struct {
	Expected *appsv1beta1.Condition
}

func (matcher *conditionMatcher) Match(actual interface{}) (success bool, err error) {
	actualCondition, ok := actual.(*appsv1beta1.Condition)
	if !ok {
		return false, errors.New("value should be a condition")
	}

	return hasSameState(actualCondition, matcher.Expected), nil
}

func (matcher *conditionMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to have the same state of", matcher.Expected)
}

func (matcher *conditionMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to have the same state of", matcher.Expected)
}
