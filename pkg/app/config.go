// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"

	"github.com/superproj/zero/pkg/log"
)

const configFlagName = "config"

var cfgFile string

// AddConfigFlag adds flags for a specific server to the specified FlagSet object.
// It also sets a passed functions to read values from configuration file into viper
// when each cobra command's Execute method is called.
func AddConfigFlag(fs *pflag.FlagSet, name string, watch bool) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.ReplaceAll(strings.ToUpper(name), "-", "_"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(".")

			if names := strings.Split(name, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homedir.HomeDir(), "."+names[0]))
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}

			viper.SetConfigName(name)
		}

		if err := viper.ReadInConfig(); err != nil {
			log.Debugw("Failed to read configuration file", "file", cfgFile, "err", err)
		}

		if watch {
			viper.WatchConfig()
			viper.OnConfigChange(func(e fsnotify.Event) {
				log.Debugw("Config file changed", "name", e.Name)
			})
		}
	})
}

func PrintConfig() {
	for _, key := range viper.AllKeys() {
		log.Debugw(fmt.Sprintf("CFG: %s=%v", key, viper.Get(key)))
	}
}

func init() {
	pflag.StringVarP(&cfgFile, configFlagName, "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}
