// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package validation

import (
	"crypto/tls"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	netutils "k8s.io/utils/net"

	"github.com/superproj/zero/internal/pkg/known"
)

const (
	DNSName              string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	SkipVerifyAnnotation        = "zero.io/skip-verify"
)

var rxDNSName = regexp.MustCompile(DNSName)

// IsValiadURL tests that https://host:port is reachble in timeout.
func IsValiadURL(url string, timeout time.Duration) error {
	client := &http.Client{
		// disabel redirect func for import clusternet proxy cluster case
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func IsValidDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		return false
	}
	return !IsValidIP(str) && rxDNSName.MatchString(str)
}

func IsValidIP(str string) bool {
	return net.ParseIP(str) != nil
}

/*
func ValidateAccountLabelsUpdate(newMeta, oldMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	//immutableLabels := []string{known.LabelAppID, known.LabelUIN, known.LabelRegion}
	immutableLabels := []string{known.LabelAppID}

	for _, label := range immutableLabels {
		value, ok := oldMeta.Labels[label]
		if ok && newMeta.Labels[label] != value {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("labels"), newMeta.Labels[label], "update not allowed"))
		}
	}

	return allErrs
}
*/

func SkipVerify(obj metav1.Object) bool {
	if !obj.GetDeletionTimestamp().IsZero() {
		return true
	}

	annotations := obj.GetAnnotations()
	if annotations != nil {
		if verify, ok := annotations[known.SkipVerifyAnnotation]; ok && verify == "true" {
			return true
		}
	}

	return false
}

func ValidateHostPort(input string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	hostIP, port, err := net.SplitHostPort(input)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, input, "must be IP:port"))
		return allErrs
	}

	if ip := netutils.ParseIPSloppy(hostIP); ip == nil {
		allErrs = append(allErrs, field.Invalid(fldPath, hostIP, "must be a valid IP"))
	}

	if p, err := strconv.Atoi(port); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, port, "must be a valid port"))
	} else if p < 1 || p > 65535 {
		allErrs = append(allErrs, field.Invalid(fldPath, port, "must be a valid port"))
	}

	return allErrs
}

func IsAdminUser(userID string) bool {
	return userID == known.AdminUserID
}
