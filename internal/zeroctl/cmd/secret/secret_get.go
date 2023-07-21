// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package secret

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	cmdutil "github.com/superproj/zero/internal/zeroctl/cmd/util"
	"github.com/superproj/zero/internal/zeroctl/util/templates"
	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
	"github.com/superproj/zero/pkg/cli/genericclioptions"
)

const (
	getUsageStr = "get SECRET_NAME"
)

// GetOptions is an options struct to support get subcommands.
type GetOptions struct {
	Name string

	GetSecretRequest *v1.GetSecretRequest
	client           v1.UserCenterClient

	genericclioptions.IOStreams
}

var (
	getExample = templates.Examples(`
		# Get a specified secret information
		zeroctl secret get foo`)

	getUsageErrStr = fmt.Sprintf("expected '%s'.\nSECRET_NAME is required arguments for the get command", getUsageStr)
)

// NewGetOptions returns an initialized GetOptions instance.
func NewGetOptions(ioStreams genericclioptions.IOStreams) *GetOptions {
	return &GetOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdGet returns new initialized instance of get sub command.
func NewCmdGet(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "get SECRET_NAME",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Display a secret resource",
		TraverseChildren:      true,
		Long:                  "Display a secret resource.",
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
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, getUsageErrStr)
	}

	o.GetSecretRequest = &v1.GetSecretRequest{Name: args[0]}
	o.client = f.UserCenterClient()

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *GetOptions) Validate(cmd *cobra.Command, args []string) error {
	return o.GetSecretRequest.Validate()
}

// Run executes a get subcommand using the specified options.
func (o *GetOptions) Run(f cmdutil.Factory, args []string) error {
	secret, err := o.client.GetSecret(f.MustWithToken(context.Background()), o.GetSecretRequest)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(o.Out)
	data := [][]string{
		{
			secret.Name,
			strconv.FormatUint(uint64(secret.Status), 10),
			secret.SecretID,
			secret.SecretKey,
			time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05"),
			secret.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
		},
	}

	table = setHeader(table)
	table = cmdutil.TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}
