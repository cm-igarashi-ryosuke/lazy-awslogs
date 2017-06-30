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

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	Aws "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/aws"
	Config "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/config"
	"github.com/spf13/cobra"
)

type reloadFlag struct {
	group string
}

var _reloadFlag reloadFlag

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Get groups and streams and update config file",
	Long:  `Get groups and streams and update config file. cache is used to the bash-completion.`,
	Run:   reloadRun,
}

func reloadRun(cmd *cobra.Command, args []string) {
	config := Config.Load()
	cwlogsClient := Aws.NewCloudWatchLogsClient(rootFlag.GetCloudWatchLogsClientParams())

	// This func returns a []Stream by a specified log group name.
	getStreams := func(logGroupName string) (streams []Config.Stream, err error) {
		err = cwlogsClient.GetLogStreams(logGroupName, func(out *cloudwatchlogs.DescribeLogStreamsOutput) {
			for _, logStream := range out.LogStreams {
				streams = append(streams, Config.Stream{
					Name: *logStream.LogStreamName,
				})
			}
		})
		return
	}

	// If a log group name is not specified, get all groups and streams.
	if _reloadFlag.group == "" {
		groups := []Config.Group{}

		err := cwlogsClient.GetLogGroups(func(out *cloudwatchlogs.DescribeLogGroupsOutput) {
			for _, logGroup := range out.LogGroups {
				streams, err := getStreams(*logGroup.LogGroupName)

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				groups = append(groups, Config.Group{
					Name:    *logGroup.LogGroupName,
					Streams: streams,
				})
			}
		})

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		config.UpdateGroups(groups)
	} else {
		streams, err := getStreams(_reloadFlag.group)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		config.UpdateStreams(_reloadFlag.group, streams)
	}
}

func init() {
	RootCmd.AddCommand(reloadCmd)

	reloadCmd.Flags().StringVarP(&_reloadFlag.group, "group", "g", "", "Specify a log group name")
}
