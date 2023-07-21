// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"

	apiservapp "github.com/superproj/zero/cmd/zero-apiserver/app"
	ctrlmgrapp "github.com/superproj/zero/cmd/zero-controller-manager/app"
	demoapp "github.com/superproj/zero/cmd/zero-demo/app"
	gwapp "github.com/superproj/zero/cmd/zero-gateway/app"
	minerctrlapp "github.com/superproj/zero/cmd/zero-miner-controller/app"
	minersetctrlapp "github.com/superproj/zero/cmd/zero-minerset-controller/app"
	watchapp "github.com/superproj/zero/cmd/zero-nightwatch/app"
	pumpapp "github.com/superproj/zero/cmd/zero-pump/app"
	toyblcapp "github.com/superproj/zero/cmd/zero-toyblc/app"
	usercenterapp "github.com/superproj/zero/cmd/zero-usercenter/app"
	zeroctlcmd "github.com/superproj/zero/internal/zeroctl/cmd"
	genutil "github.com/superproj/zero/pkg/util/gen"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path, module := "", ""
	if len(os.Args) == 3 {
		path = os.Args[1]
		module = os.Args[2]
	} else {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory] [module] \n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutil.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	switch module {
	case "zero-demo":
		// generate docs for zero-demo
		demo := demoapp.NewApp("zero-demo").Command()
		_ = doc.GenMarkdownTree(demo, outDir)
	case "zero-usercenter":
		// generate docs for zero-usercenter
		usercenter := usercenterapp.NewApp("zero-usercenter").Command()
		_ = doc.GenMarkdownTree(usercenter, outDir)
	case "zero-apiserver":
		// generate docs for zero-apiserver
		apiserver := apiservapp.NewAPIServerCommand()
		_ = doc.GenMarkdownTree(apiserver, outDir)
	case "zero-gateway":
		// generate docs for zero-gateway
		gwserver := gwapp.NewApp("zero-gateway").Command()
		_ = doc.GenMarkdownTree(gwserver, outDir)
	case "zero-nightwatch":
		// generate docs for zero-nightwatch
		nw := watchapp.NewApp("zero-nightwatch").Command()
		_ = doc.GenMarkdownTree(nw, outDir)
	case "zero-pump":
		// generate docs for zero-pump
		pump := pumpapp.NewApp("zero-pump").Command()
		_ = doc.GenMarkdownTree(pump, outDir)
	case "zero-toyblc":
		// generate docs for zero-toyblc
		toyblc := toyblcapp.NewApp("zero-toyblc").Command()
		_ = doc.GenMarkdownTree(toyblc, outDir)
	case "zero-controller-manager":
		// generate docs for zero-controller-manager
		ctrlmgr := ctrlmgrapp.NewControllerManagerCommand()
		_ = doc.GenMarkdownTree(ctrlmgr, outDir)
	case "zero-minerset-controller":
		// generate docs for zero-minerset-controller
		minersetctrl := minersetctrlapp.NewControllerCommand()
		_ = doc.GenMarkdownTree(minersetctrl, outDir)
	case "zero-miner-controller":
		// generate docs for zero-miner-controller
		minerctrl := minerctrlapp.NewControllerCommand()
		_ = doc.GenMarkdownTree(minerctrl, outDir)
	case "zeroctl":
		// generate docs for zeroctl
		zeroctl := zeroctlcmd.NewDefaultZeroCtlCommand()
		_ = doc.GenMarkdownTree(zeroctl, outDir)
	default:
		fmt.Fprintf(os.Stderr, "Module %s is not supported", module)
		os.Exit(1)
	}
}
