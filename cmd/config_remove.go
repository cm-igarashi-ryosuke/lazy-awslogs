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

// FIXME: Add help message for param. Currently param is not flag.

// removeCmd represents the remove command
var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an Environment",
	Long:  `Remove an Environment.`,
	Run:   configRemoveRun,
}

func configRemoveRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Printf("Args expects to have 1, but it have %d\n", len(args))
		os.Exit(0)
	}
	config := Config.Load()
	if err := config.RemoveEnvironment(args[0]); err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}

func init() {
	configCmd.AddCommand(configRemoveCmd)
}
