// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package cmd create a root cobra command and add subcommands to it.
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/superproj/zero/internal/zeroctl/cmd/color"
	"github.com/superproj/zero/internal/zeroctl/cmd/completion"
	"github.com/superproj/zero/internal/zeroctl/cmd/info"
	"github.com/superproj/zero/internal/zeroctl/cmd/jwt"
	"github.com/superproj/zero/internal/zeroctl/cmd/new"
	"github.com/superproj/zero/internal/zeroctl/cmd/options"
	"github.com/superproj/zero/internal/zeroctl/cmd/user"
	cmdutil "github.com/superproj/zero/internal/zeroctl/cmd/util"
	"github.com/superproj/zero/internal/zeroctl/cmd/validate"
	"github.com/superproj/zero/internal/zeroctl/cmd/version"

	// "github.com/superproj/zero/internal/zeroctl/plugin".
	"github.com/superproj/zero/internal/zeroctl/cmd/secret"
	clioptions "github.com/superproj/zero/internal/zeroctl/util/options"
	"github.com/superproj/zero/internal/zeroctl/util/templates"
	"github.com/superproj/zero/internal/zeroctl/util/term"
	"github.com/superproj/zero/pkg/cli/genericclioptions"
)

// NewDefaultZeroCtlCommand creates the `zeroctl` command with default arguments.
func NewDefaultZeroCtlCommand() *cobra.Command {
	return NewZeroCtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

// NewZeroCtlCommand returns new initialized instance of 'zeroctl' root command.
func NewZeroCtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	warningHandler := rest.NewWarningWriter(err, rest.WarningWriterOptions{Deduplicate: true, Color: term.AllowsColorOutput(err)})
	warningsAsErrors := false
	zeroCliOptions := clioptions.NewOptions()
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "zeroctl",
		Short: "zeroctl controls the zero cloud platform",
		Long: templates.LongDesc(`
		zeroctl controls the zero cloud platform, is the client side tool for zero cloud platform.

		Find more information at:
			https://github.com/superproj/zero/blob/master/docs/guide/en-US/cmd/zeroctl/zeroctl.md`),
		Run: runHelp,
		// Hook before and after Run initialize and write profiles to disk,
		// respectively.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			rest.SetDefaultWarningHandler(warningHandler)

			if cmd.Name() == cobra.ShellCompRequestCmd {
				// This is the __complete or __completeNoDesc command which
				// indicates shell completion has been requested.
				// plugin.SetupPluginCompletion(cmd, args)
			}

			zeroCliOptions.Complete()

			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			if err := flushProfiling(); err != nil {
				return err
			}
			if warningsAsErrors {
				count := warningHandler.WarningCount()
				switch count {
				case 0:
					// no warnings
				case 1:
					return fmt.Errorf("%d warning received", count)
				default:
					return fmt.Errorf("%d warnings received", count)
				}
			}
			return nil
		},
	}

	// From this point and forward we get warnings on flags that contain "_" separators
	// when adding them with hyphen instead of the original name.
	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	flags := cmds.PersistentFlags()

	addProfilingFlags(flags)

	flags.BoolVar(&warningsAsErrors, "warnings-as-errors", warningsAsErrors, "Treat warnings received from the server as errors and exit with a non-zero exit code")

	zeroCliOptions.AddFlags(flags)

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	_ = viper.BindPFlags(cmds.PersistentFlags())
	cobra.OnInitialize(func() {
		initConfig(viper.GetString(clioptions.FlagZeroConfig))
	})
	cmds.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	f := cmdutil.NewFactory(zeroCliOptions)

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands:",
			Commands: []*cobra.Command{
				info.NewCmdInfo(f, ioStreams),
				color.NewCmdColor(f, ioStreams),
				new.NewCmdNew(f, ioStreams),
				jwt.NewCmdJWT(f, ioStreams),
			},
		},
		{
			Message: "UserCenter Commands:",
			Commands: []*cobra.Command{
				user.NewCmdUser(f, ioStreams),
				secret.NewCmdSecret(f, ioStreams),
			},
		},
		{
			Message:  "Gateway Commands:",
			Commands: []*cobra.Command{},
		},
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				validate.NewCmdValidate(f, ioStreams),
			},
		},
		{
			Message: "Settings Commands:",
			Commands: []*cobra.Command{
				// set.NewCmdSet(f, ioStreams),
				completion.NewCmdCompletion(ioStreams.Out, ""),
			},
		},
	}
	groups.Add(cmds)

	filters := []string{"options"}

	/*
		// Hide the "alpha" subcommand if there are no alpha commands in this build.
		alpha := NewCmdAlpha(f, ioStreams)
		if !alpha.HasSubCommands() {
			filters = append(filters, alpha.Name())
		}
	*/

	templates.ActsAsRootCommand(cmds, filters, groups...)

	// cmds.AddCommand(alpha)
	// cmds.AddCommand(plugin.NewCmdPlugin(ioStreams))
	cmds.AddCommand(version.NewCmdVersion(f, ioStreams))
	cmds.AddCommand(options.NewCmdOptions(ioStreams.Out))

	// Stop warning about normalization of flags. That makes it possible to
	// add the klog flags later.
	cmds.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)
	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
