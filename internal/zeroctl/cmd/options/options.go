// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package options print a list of global command-line options (applies to all commands).
package options

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/superproj/zero/internal/zeroctl/util/templates"
)

var optionsExample = templates.Examples(`
		# Print flags inherited by all commands
		zeroctl options`)

// NewCmdOptions implements the options command.
func NewCmdOptions(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "options",
		Short:   "Print the list of flags inherited by all commands",
		Long:    "Print the list of flags inherited by all commands",
		Example: optionsExample,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	// The `options` command needs write its output to the `out` stream
	// (typically stdout). Without calling SetOutput here, the Usage()
	// function call will fall back to stderr.
	cmd.SetOutput(out)
	cmd.SetErr(out)

	templates.UseOptionsTemplates(cmd)

	return cmd
}
