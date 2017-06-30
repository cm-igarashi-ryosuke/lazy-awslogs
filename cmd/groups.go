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

	Aws "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/aws"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	Config "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/config"
	"github.com/spf13/cobra"
)

type groupFlag struct {
	cache   bool
	verbose bool
}

var _groupFlag groupFlag

// groupsCmd represents the groups command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List groups",
	Long:  `List groups.`,
	Run:   groupsRun,
}

func groupsRun(cmd *cobra.Command, args []string) {
	if rootFlag.verbose {
		fmt.Printf("profile=%s\n", rootFlag.profile)
	}

	if _groupFlag.cache {
		config := Config.Load()
		groups := config.ReadGroups()
		for _, group := range groups {
			fmt.Println(group.Name)
		}
	} else {
		cwlogsClient := Aws.NewCloudWatchLogsClient(rootFlag.GetCloudWatchLogsClientParams())

		err := cwlogsClient.GetLogGroups(func(out *cloudwatchlogs.DescribeLogGroupsOutput) {
			for _, logGroup := range out.LogGroups {
				fmt.Println(*logGroup.LogGroupName)
			}
		})

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func init() {
	RootCmd.AddCommand(groupsCmd)

	groupsCmd.Flags().BoolVarP(&_groupFlag.cache, "cache", "c", false, "Use cache in config file")
}
