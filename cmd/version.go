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
	"runtime"

	"github.com/spf13/cobra"
)

var appVersion string = "1.0.0" // set value with link

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  `Display version information and exit.`,
	Run:   versionRun,
}

func versionRun(cmd *cobra.Command, args []string) {
	app := "lazy-awslogs"
	fmt.Printf("%s %s (%s, %s/%s)\n", app, appVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
