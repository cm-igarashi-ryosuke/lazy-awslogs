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

	Config "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Environments",
	Long:  `List Environments.`,
	Run:   configListRun,
}

func configListRun(cmd *cobra.Command, args []string) {
	config := Config.Load()
	for _, env := range config.Environments {
		prefix := "  "
		if env.Name == config.Current {
			prefix = "* "
		}
		fmt.Println(prefix + env.Name)
	}
}

func init() {
	configCmd.AddCommand(configListCmd)
}
