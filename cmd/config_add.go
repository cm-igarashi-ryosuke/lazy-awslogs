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
	"fmt"
	"os"

	Config "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/config"
	"github.com/spf13/cobra"
)

type configAddFlag struct {
	name    string
	profile string
	region  string
}

var _configAddFlag configAddFlag

// addCmd represents the add command
var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new Environment",
	Long:  `Add a new Environment.`,
	Run:   configAddRun,
}

func configAddRun(cmd *cobra.Command, args []string) {
	if _configAddFlag.name == "" {
		fmt.Println("Require --config-name option.")
		os.Exit(0)
	}
	config := Config.Load()
	newEnvironment := Config.Environment{
		Name:    _configAddFlag.name,
		Profile: _configAddFlag.profile,
		Region:  _configAddFlag.region,
	}
	if err := config.AddEnvironment(newEnvironment); err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}

func init() {
	configCmd.AddCommand(configAddCmd)
	configAddCmd.Flags().StringVarP(&_configAddFlag.name, "config-name", "", "", "[required] New environment's name")
	configAddCmd.Flags().StringVarP(&_configAddFlag.profile, "config-profile", "", "", "Set AWS credentials profile")
	configAddCmd.Flags().StringVarP(&_configAddFlag.region, "config-region", "", "", "Set default AWS Region")
}
