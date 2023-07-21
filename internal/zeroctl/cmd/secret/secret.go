// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package secret provides functions to manage secrets on zero platform.
package secret

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	cmdutil "github.com/superproj/zero/internal/zeroctl/cmd/util"
	"github.com/superproj/zero/internal/zeroctl/util/templates"
	"github.com/superproj/zero/pkg/cli/genericclioptions"
)

var secretLong = templates.LongDesc(`
	Secret management commands.

	This commands allow you to manage your secret on zero platform.`)

// NewCmdSecret returns new initialized instance of 'secret' sub command.
func NewCmdSecret(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "secret SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage secrets on zero platform",
		Long:                  secretLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.ErrOut),
	}

	cmd.AddCommand(NewCmdCreate(f, ioStreams))
	cmd.AddCommand(NewCmdGet(f, ioStreams))
	cmd.AddCommand(NewCmdList(f, ioStreams))
	cmd.AddCommand(NewCmdDelete(f, ioStreams))
	cmd.AddCommand(NewCmdUpdate(f, ioStreams))

	return cmd
}

// setHeader set headers for secret commands.
func setHeader(table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{"Name", "Status", "SecretID", "SecretKey", "Expires", "Created"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.FgGreenColor})

	return table
}
