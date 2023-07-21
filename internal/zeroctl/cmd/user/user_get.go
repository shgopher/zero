// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package user

import (
	"context"
	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	cmdutil "github.com/superproj/zero/internal/zeroctl/cmd/util"
	"github.com/superproj/zero/internal/zeroctl/util/templates"
	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
	"github.com/superproj/zero/pkg/cli/genericclioptions"
)

const (
	getUsageStr = "get USERNAME"
)

// GetOptions is an options struct to support get subcommands.
type GetOptions struct {
	Name string

	GetUserRequest *v1.GetUserRequest
	client         v1.UserCenterClient

	genericclioptions.IOStreams
}

var (
	getExample = templates.Examples(`
		# Get user foo detail information
		zeroctl user get foo

		# Get current login user detail information
		zeroctl user info`)

	getUsageErrStr = fmt.Sprintf("expected '%s'.\nUSERNAME is required arguments for the get command", getUsageStr)
)

// NewGetOptions returns an initialized GetOptions instance.
func NewGetOptions(ioStreams genericclioptions.IOStreams) *GetOptions {
	return &GetOptions{
		GetUserRequest: &v1.GetUserRequest{},
		IOStreams:      ioStreams,
	}
}

// NewCmdGet returns new initialized instance of get sub command.
func NewCmdGet(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   getUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{"info"},
		Short:                 "Display a user resource.",
		TraverseChildren:      true,
		Long:                  `Display a user resource.`,
		Example:               getExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(f, args))
		},
		SuggestFor: []string{},
	}

	return cmd
}

// Complete completes all the required options.
func (o *GetOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	o.GetUserRequest.Username = f.GetOptions().UserOptions.Username
	if len(args) != 0 {
		o.GetUserRequest.Username = args[0]
	}

	o.client = f.UserCenterClient()

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *GetOptions) Validate(cmd *cobra.Command, args []string) error {
	if cmd.CalledAs() == "info" && len(args) > 0 {
		return fmt.Errorf("%q does not take any arguments, got %q", cmd.CalledAs(), args)
	}

	if cmd.CalledAs() == "get" && len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, getUsageErrStr)
	}

	return o.GetUserRequest.Validate()
}

// Run executes a get subcommand using the specified options.
func (o *GetOptions) Run(f cmdutil.Factory, args []string) error {
	user, err := o.client.GetUser(f.MustWithToken(context.Background()), o.GetUserRequest)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(o.Out)

	data := [][]string{
		{
			user.Username,
			user.UserID,
			user.Nickname,
			user.Email,
			user.Phone,
			user.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
			user.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"),
		},
	}

	table = setHeader(table)
	table = cmdutil.TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}
