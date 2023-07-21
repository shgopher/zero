// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:funlen,gocritic
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	mangen "github.com/cpuguy83/go-md2man/v2/md2man"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	apiservapp "github.com/superproj/zero/cmd/zero-apiserver/app"
	ctrlmgrapp "github.com/superproj/zero/cmd/zero-controller-manager/app"
	demoapp "github.com/superproj/zero/cmd/zero-demo/app"
	gwapp "github.com/superproj/zero/cmd/zero-gateway/app"
	minerctrlapp "github.com/superproj/zero/cmd/zero-miner-controller/app"
	minersetctrlapp "github.com/superproj/zero/cmd/zero-minerset-controller/app"
	nwapp "github.com/superproj/zero/cmd/zero-nightwatch/app"
	pumpapp "github.com/superproj/zero/cmd/zero-pump/app"
	toyblcapp "github.com/superproj/zero/cmd/zero-toyblc/app"
	usercenterapp "github.com/superproj/zero/cmd/zero-usercenter/app"
	zeroctlcmd "github.com/superproj/zero/internal/zeroctl/cmd"
	genutil "github.com/superproj/zero/pkg/util/gen"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/man/man1"
	module := ""
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

	// Set environment variables used by command so the output is consistent,
	// regardless of where we run.
	_ = os.Setenv("HOME", "/home/username")

	switch module {
	case "zero-demo":
		// generate manpage for zero-demo
		demo := demoapp.NewApp("zero-demo").Command()
		genMarkdown(demo, "", outDir)
		for _, c := range demo.Commands() {
			genMarkdown(c, "zero-demo", outDir)
		}
	case "zero-usercenter":
		// generate manpage for zero-usercenter
		usercenter := usercenterapp.NewApp("zero-usercenter").Command()
		genMarkdown(usercenter, "", outDir)
		for _, c := range usercenter.Commands() {
			genMarkdown(c, "zero-usercenter", outDir)
		}
	case "zero-apiserver":
		// generate manpage for zero-apiserver
		apiserver := apiservapp.NewAPIServerCommand()
		genMarkdown(apiserver, "", outDir)
		for _, c := range apiserver.Commands() {
			genMarkdown(c, "zero-apiserver", outDir)
		}
	case "zero-gateway":
		// generate manpage for zero-gateway
		gwserver := gwapp.NewApp("zero-gateway").Command()
		genMarkdown(gwserver, "", outDir)
		for _, c := range gwserver.Commands() {
			genMarkdown(c, "zero-gateway", outDir)
		}
	case "zero-nightwatch":
		// generate manpage for zero-nightwatch
		nw := nwapp.NewApp("zero-nightwatch").Command()
		genMarkdown(nw, "", outDir)
		for _, c := range nw.Commands() {
			genMarkdown(c, "zero-nightwatch", outDir)
		}
	case "zero-pump":
		// generate manpage for zero-pump
		pump := pumpapp.NewApp("zero-pump").Command()
		genMarkdown(pump, "", outDir)
		for _, c := range pump.Commands() {
			genMarkdown(c, "zero-pump", outDir)
		}
	case "zero-toyblc":
		// generate manpage for zero-toyblc
		toyblc := toyblcapp.NewApp("zero-toyblc").Command()
		genMarkdown(toyblc, "", outDir)
		for _, c := range toyblc.Commands() {
			genMarkdown(c, "zero-toyblc", outDir)
		}
	case "zero-controller-manager":
		// generate manpage for zero-controller-manager
		ctrlmgr := ctrlmgrapp.NewControllerManagerCommand()
		genMarkdown(ctrlmgr, "", outDir)
		for _, c := range ctrlmgr.Commands() {
			genMarkdown(c, "zero-controller-manager", outDir)
		}
	case "zero-minerset-controller":
		// generate manpage for zero-minerset-controller
		minersetctrl := minersetctrlapp.NewControllerCommand()
		genMarkdown(minersetctrl, "", outDir)
		for _, c := range minersetctrl.Commands() {
			genMarkdown(c, "zero-minerset-controller", outDir)
		}
	case "zero-miner-controller":
		// generate manpage for zero-miner-controller
		minerctrl := minerctrlapp.NewControllerCommand()
		genMarkdown(minerctrl, "", outDir)
		for _, c := range minerctrl.Commands() {
			genMarkdown(c, "zero-miner-controller", outDir)
		}
	case "zeroctl":
		// generate manpage for zeroctl
		// TODO os.Stdin should really be something like ioutil.Discard, but a Reader
		zeroctl := zeroctlcmd.NewDefaultZeroCtlCommand()
		genMarkdown(zeroctl, "", outDir)
		for _, c := range zeroctl.Commands() {
			genMarkdown(c, "zeroctl", outDir)
		}
	default:
		fmt.Fprintf(os.Stderr, "Module %s is not supported", module)
		os.Exit(1)
	}
}

func preamble(out *bytes.Buffer, name, short, long string) {
	out.WriteString(`% Zero(1) zero User Manuals
% Eric Paris
% Jan 2015
# NAME
`)
	fmt.Fprintf(out, "%s \\- %s\n\n", name, short)
	fmt.Fprintf(out, "# SYNOPSIS\n")
	fmt.Fprintf(out, "**%s** [OPTIONS]\n\n", name)
	fmt.Fprintf(out, "# DESCRIPTION\n")
	fmt.Fprintf(out, "%s\n\n", long)
}

func printFlags(out io.Writer, flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		format := "**--%s**=%s\n\t%s\n\n"
		if flag.Value.Type() == "string" {
			// put quotes on the value
			format = "**--%s**=%q\n\t%s\n\n"
		}

		// Todo, when we mark a shorthand is deprecated, but specify an empty message.
		// The flag.ShorthandDeprecated is empty as the shorthand is deprecated.
		// Using len(flag.ShorthandDeprecated) > 0 can't handle this, others are ok.
		if !(len(flag.ShorthandDeprecated) > 0) && len(flag.Shorthand) > 0 {
			format = "**-%s**, " + format
			fmt.Fprintf(out, format, flag.Shorthand, flag.Name, flag.DefValue, flag.Usage)
		} else {
			fmt.Fprintf(out, format, flag.Name, flag.DefValue, flag.Usage)
		}
	})
}

func printOptions(out io.Writer, command *cobra.Command) {
	flags := command.NonInheritedFlags()
	if flags.HasFlags() {
		fmt.Fprintf(out, "# OPTIONS\n")
		printFlags(out, flags)
		fmt.Fprintf(out, "\n")
	}
	flags = command.InheritedFlags()
	if flags.HasFlags() {
		fmt.Fprintf(out, "# OPTIONS INHERITED FROM PARENT COMMANDS\n")
		printFlags(out, flags)
		fmt.Fprintf(out, "\n")
	}
}

func genMarkdown(command *cobra.Command, parent, docsDir string) {
	dparent := strings.ReplaceAll(parent, " ", "-")
	name := command.Name()

	dname := name
	if len(parent) > 0 {
		dname = dparent + "-" + name
		name = parent + " " + name
	}

	out := new(bytes.Buffer)

	short, long := command.Short, command.Long
	if len(long) == 0 {
		long = short
	}

	preamble(out, name, short, long)
	printOptions(out, command)

	if len(command.Example) > 0 {
		fmt.Fprintf(out, "# EXAMPLE\n")
		fmt.Fprintf(out, "```\n%s\n```\n", command.Example)
	}

	if len(command.Commands()) > 0 || len(parent) > 0 {
		fmt.Fprintf(out, "# SEE ALSO\n")

		if len(parent) > 0 {
			fmt.Fprintf(out, "**%s(1)**, ", dparent)
		}

		for _, c := range command.Commands() {
			fmt.Fprintf(out, "**%s-%s(1)**, ", dname, c.Name())
			genMarkdown(c, name, docsDir)
		}

		fmt.Fprintf(out, "\n")
	}

	out.WriteString(`
# HISTORY
January 2015, Originally compiled by Eric Paris (eparis at redhat dot com) based on the superproj source material, but hopefully they have been automatically generated since!
`)

	final := mangen.Render(out.Bytes())

	filename := docsDir + dname + ".1"

	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	_, err = outFile.Write(final)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
