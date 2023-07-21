// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/superproj/zero/internal/zeroctl/cmd"
	"github.com/superproj/zero/pkg/util/gen"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/"
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else if len(os.Args) > 2 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [output directory]\n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := gen.OutDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used by zeroctl so the output is consistent,
	// regardless of where we run.
	_ = os.Setenv("HOME", "/home/username")
	// TODO os.Stdin should really be something like ioutil.Discard, but a Reader
	zeroctl := cmd.NewZeroCtlCommand(os.Stdin, ioutil.Discard, ioutil.Discard)
	_ = doc.GenMarkdownTree(zeroctl, outDir)
}
