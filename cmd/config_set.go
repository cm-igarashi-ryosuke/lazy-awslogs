// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	Config "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/config"
	"github.com/spf13/cobra"
)

type configSetFlag struct {
	profile string
	region  string
	group   string
	stream  string
}

var _configSetFlag configSetFlag

// setCmd represents the set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set current config settings",
	Long:  `Set current config settings.`,
	Run:   configSetRun,
}

func configSetRun(cmd *cobra.Command, args []string) {
	config := Config.Load()
	if _configSetFlag.profile != "" {
		config.SetCurrentEnvironmentParam(Config.Profile, _configSetFlag.profile)
	}
	if _configSetFlag.region != "" {
		config.SetCurrentEnvironmentParam(Config.Region, _configSetFlag.region)
	}
	if _configSetFlag.group != "" {
		config.SetCurrentEnvironmentParam(Config.DefaultGroup, _configSetFlag.group)
	}
	if _configSetFlag.stream != "" {
		config.SetCurrentEnvironmentParam(Config.DefaultStream, _configSetFlag.stream)
	}
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configSetCmd.Flags().StringVarP(&_configSetFlag.profile, "config-profile", "", "", "Set default AWS credentials profile")
	configSetCmd.Flags().StringVarP(&_configSetFlag.region, "config-region", "", "", "Set default AWS Region")
	configSetCmd.Flags().StringVarP(&_configSetFlag.group, "config-group", "", "", "Set default log group name")
	configSetCmd.Flags().StringVarP(&_configSetFlag.stream, "config-stream", "", "", "Set default log stream name")
}
