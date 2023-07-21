// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:gocritic
package main

import (
	"fmt"
	"io"
	"os"

	kruntime "k8s.io/apimachinery/pkg/runtime"

	flag "github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	functionDest = flag.StringP("func-dest", "f", "-", "Output for swagger functions; '-' means stdout (default)")
	typeSrc      = flag.StringP("type-src", "s", "", "From where we are going to read the types")
	verify       = flag.BoolP("verify", "v", false, "Verifies if the given type-src file has documentation for every type")
)

func main() {
	flag.Parse()

	if *typeSrc == "" {
		klog.Fatalf("Please define -s flag as it is the source file")
	}

	var funcOut io.Writer
	if *functionDest == "-" {
		funcOut = os.Stdout
	} else {
		file, err := os.Create(*functionDest)
		if err != nil {
			klog.Fatalf("Couldn't open %v: %v", *functionDest, err)
		}
		defer file.Close()
		funcOut = file
	}

	docsForTypes := kruntime.ParseDocumentationFrom(*typeSrc)

	if *verify {
		rc, err := kruntime.VerifySwaggerDocsExist(docsForTypes, funcOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in verification process: %s\n", err)
		}
		os.Exit(rc)
	}

	if len(docsForTypes) > 0 {
		if err := kruntime.WriteSwaggerDocFunc(docsForTypes, funcOut); err != nil {
			fmt.Fprintf(os.Stderr, "Error when writing swagger documentation functions: %s\n", err)
			os.Exit(-1)
		}
	}
}
