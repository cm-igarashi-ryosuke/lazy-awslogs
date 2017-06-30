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

type streamsFlag struct {
	group string
	cache bool
}

var _streamsFlag streamsFlag

// streamsCmd represents the streams command
var streamsCmd = &cobra.Command{
	Use:   "streams",
	Short: "List streams",
	Long:  `List streams.`,
	Run:   streamsRun,
}

func streamsRun(cmd *cobra.Command, args []string) {
	if rootFlag.verbose {
		fmt.Printf("profile=%s\n", rootFlag.profile)
	}

	if _streamsFlag.cache {
		config := Config.Load()
		streams := config.ReadStreams(_streamsFlag.group)
		for _, stream := range streams {
			fmt.Println(stream.Name)
		}
	} else {
		cwlogsClient := Aws.NewCloudWatchLogsClient(rootFlag.GetCloudWatchLogsClientParams())

		err := cwlogsClient.GetLogStreams(_streamsFlag.group, func(out *cloudwatchlogs.DescribeLogStreamsOutput) {
			for _, logStream := range out.LogStreams {
				fmt.Println(*logStream.LogStreamName)
			}
		})

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func init() {
	RootCmd.AddCommand(streamsCmd)

	streamsCmd.Flags().StringVarP(&_streamsFlag.group, "group", "g", "", "[required] Specify a log group name")
	streamsCmd.Flags().BoolVarP(&_streamsFlag.cache, "cache", "c", false, "Use cache in config file")
}
