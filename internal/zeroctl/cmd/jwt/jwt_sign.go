// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"

	cmdutil "github.com/superproj/zero/internal/zeroctl/cmd/util"
	"github.com/superproj/zero/internal/zeroctl/util/templates"
	"github.com/superproj/zero/pkg/cli/genericclioptions"
)

const (
	signUsageStr = "sign SECRETID SECRETKEY"
)

// ErrSigningMethod defines invalid signing method error.
var ErrSigningMethod = errors.New("invalid signing method")

// SignOptions is an options struct to support sign subcommands.
type SignOptions struct {
	Timeout   time.Duration
	NotBefore time.Duration
	Algorithm string
	Issuer    string
	Header    ArgList

	genericclioptions.IOStreams
}

var (
	signExample = templates.Examples(`
		# Sign a token with secretID and secretKey
		zeroctl sign tgydj8d9EQSnFqKf iBdEdFNBLN1nR3fV

		# Sign a token with expires and sign method
		zeroctl sign tgydj8d9EQSnFqKf iBdEdFNBLN1nR3fV --timeout=2h --algorithm=HS256`)

	signUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nSECRETID and SECRETKEY are required arguments for the sign command",
		signUsageStr,
	)
)

// NewSignOptions returns an initialized SignOptions instance.
func NewSignOptions(ioStreams genericclioptions.IOStreams) *SignOptions {
	return &SignOptions{
		Timeout:   2 * time.Hour,
		Algorithm: "HS256",
		Issuer:    "zeroctl",
		Header:    make(ArgList),

		IOStreams: ioStreams,
	}
}

// NewCmdSign returns new initialized instance of sign sub command.
func NewCmdSign(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewSignOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   signUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Sign a jwt token with given secretID and secretKey",
		Long:                  "Sign a jwt token with given secretID and secretKey",
		TraverseChildren:      true,
		Example:               signExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return cmdutil.UsageErrorf(cmd, signUsageErrStr)
			}

			return nil
		},
	}

	// mark flag as deprecated
	cmd.Flags().DurationVar(&o.Timeout, "timeout", o.Timeout, "JWT token expires time.")
	cmd.Flags().DurationVar(&o.NotBefore, "not-before", o.NotBefore, "Identifies the time before which the JWT MUST NOT be accepted for processing.")
	cmd.Flags().StringVar(&o.Algorithm, "algorithm", o.Algorithm, "Signing algorithm - possible values are HS256, HS384, HS512.")
	cmd.Flags().StringVar(&o.Issuer, "issuer", o.Issuer, "Identifies the principal that issued the JWT.")
	cmd.Flags().Var(&o.Header, "header", "Add additional header params. may be used more than once.")

	return cmd
}

// Complete completes all the required options.
func (o *SignOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *SignOptions) Validate(cmd *cobra.Command, args []string) error {
	switch o.Algorithm {
	case "HS256", "HS384", "HS512":
	default:
		return ErrSigningMethod
	}

	return nil
}

// Run executes a sign subcommand using the specified options.
func (o *SignOptions) Run(args []string) error {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    o.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(o.Timeout)),
		NotBefore: jwt.NewNumericDate(now),
	}

	// create a new token
	token := jwt.NewWithClaims(jwt.GetSigningMethod(o.Algorithm), claims)

	// add command line headers
	if len(o.Header) > 0 {
		for k, v := range o.Header {
			token.Header[k] = v
		}
	}
	token.Header["kid"] = args[0]

	accessToken, err := token.SignedString([]byte(args[1]))
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, accessToken+"\n")

	return nil
}
