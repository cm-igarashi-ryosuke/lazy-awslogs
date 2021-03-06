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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	Aws "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/aws"
	"github.com/spf13/cobra"

	"github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/flags"
)

type GetFlags struct {
	Profile   string
	Region    string
	Log       flags.CWLogIdentifyFlags
	Pattern   string
	TimeRange flags.CWLogTimeRange
	Watch     bool
	noPrefix  bool
}

const watchInterval = 5 // seconds

// Create FilterLogEventsInput by GetFlags
func (this *GetFlags) GetCloudWatchLogsFilterLogEventsParam() cloudwatchlogs.FilterLogEventsInput {
	input := cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: &this.Log.Group,
	}
	if this.Pattern != "" {
		pattern := fmt.Sprintf("\"%s\"", this.Pattern)
		input.FilterPattern = &pattern
	}
	if this.Log.Stream != "" {
		input.LogStreamNames = []*string{&this.Log.Stream}
	}
	if startTime := this.TimeRange.StartTimeMilliseconds(); startTime != 0 {
		input.StartTime = &startTime
	}
	if endTime := this.TimeRange.EndTimeMilliseconds(); endTime != 0 {
		input.EndTime = &endTime
	}
	return input
}

// Create GetLogEventsInput by GetFlags
func (this *GetFlags) GetCloudWatchLogsGetLogEventsParam() cloudwatchlogs.GetLogEventsInput {
	input := cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &this.Log.Group,
		LogStreamName: &this.Log.Stream,
	}
	if startTime := this.TimeRange.StartTimeMilliseconds(); startTime != 0 {
		input.StartTime = &startTime
	}
	if endTime := this.TimeRange.EndTimeMilliseconds(); endTime != 0 {
		input.EndTime = &endTime
	}
	return input
}

// Print any events to stdout
func printEvents(any interface{}, appendPrefix bool) {
	if events, ok := any.([]*cloudwatchlogs.FilteredLogEvent); ok {
		for _, event := range events {
			if appendPrefix {
				fmt.Print("[" + time.Unix(*event.Timestamp/1000, 0).String() + "] ")
				fmt.Print("[" + *event.LogStreamName + "] ")
			}
			fmt.Println(strings.TrimRight(*event.Message, "\n"))
		}
	} else if events, ok := any.([]*cloudwatchlogs.OutputLogEvent); ok {
		for _, event := range events {
			if appendPrefix {
				fmt.Print("[" + time.Unix(*event.Timestamp/1000, 0).String() + "] ")
			}
			fmt.Println(strings.TrimRight(*event.Message, "\n"))
		}
	}
}

var _getFlags = &GetFlags{}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Lists log events from the specified log stream.",
	Long: `Lists log events from the specified log stream. You can list all the log events or filter using a time range.

By default, this operation returns as many log events as can fit in a response  size of 1MB (up to 10,000 log events).
If the results include tokens, there are more log events available.
You can get additional log events by specifying one of the tokens in a subsequent call.`,
	Run:    getRun,
	PreRun: preRun,
}

func preRun(cmd *cobra.Command, args []string) {
	if rootFlag.profile == "" {
		fmt.Printf("%v\n\n", fmt.Errorf("Error: Profile is not specified."))
		cmd.Usage()
		os.Exit(1)
	}

	if _getFlags.Log.Group == "" {
		fmt.Printf("%v\n\n", fmt.Errorf("Error: Group is not specified."))
		cmd.Usage()
		os.Exit(1)
	}

	if rootFlag.verbose {
		fmt.Printf("Local flags: %#v\n", _getFlags)
	}
}

func getRun(cmd *cobra.Command, args []string) {
	cwlogsClient := Aws.NewCloudWatchLogsClient(rootFlag.GetCloudWatchLogsClientParams())
	var err error

	if _getFlags.Watch {
		// Initialize startTime and endTime
		var startTime int64
		if _getFlags.TimeRange.StartTimeMilliseconds() != 0 {
			startTime = _getFlags.TimeRange.StartTimeMilliseconds()
		} else {
			startTime = time.Now().Add(-1*time.Duration(1)*time.Minute).Unix() * 1000
		}
		endTime := time.Now().Unix() * 1000

	WatchStart:
		if _getFlags.Pattern != "" || _getFlags.Log.Stream == "" {
			options := _getFlags.GetCloudWatchLogsFilterLogEventsParam()
			options.StartTime = &startTime
			options.EndTime = &endTime
			if rootFlag.verbose {
				fmt.Printf("FilterLogEvents: options=%v\n", options)
			}
			err = cwlogsClient.FilterLogEvents(&options, func(out *cloudwatchlogs.FilterLogEventsOutput) {
				printEvents(out.Events, !_getFlags.noPrefix)
				// Update startTime and endTime
				if l := len(out.Events); l > 0 {
					startTime = *out.Events[l-1].Timestamp
					startTime = startTime + 1
				}
			})
		} else {
			options := _getFlags.GetCloudWatchLogsGetLogEventsParam()
			options.StartTime = &startTime
			options.EndTime = &endTime
			if rootFlag.verbose {
				fmt.Printf("GetLogEvents: options=%v\n", options)
			}
			err = cwlogsClient.GetLogEvents(&options, func(out *cloudwatchlogs.GetLogEventsOutput) {
				printEvents(out.Events, !_getFlags.noPrefix)
				// Update startTime and endTime
				if l := len(out.Events); l > 0 {
					startTime = *out.Events[l-1].Timestamp
					startTime = startTime + 1
				}
			})
		}
		if err == nil {
			if rootFlag.verbose {
				fmt.Printf("`watch` enabled. now sleep(%d) Zzz...\n", watchInterval)
			}
			time.Sleep(watchInterval * time.Second)
			if rootFlag.verbose {
				fmt.Println("`watch` enabled. wakeup now!")
			}
			endTime = time.Now().Unix() * 1000
			goto WatchStart
		} else {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	} else {
		if _getFlags.Pattern != "" || _getFlags.Log.Stream == "" {
			options := _getFlags.GetCloudWatchLogsFilterLogEventsParam()
			if rootFlag.verbose {
				fmt.Printf("FilterLogEvents: options=%v\n", options)
			}
			err = cwlogsClient.FilterLogEvents(&options, func(out *cloudwatchlogs.FilterLogEventsOutput) {
				printEvents(out.Events, !_getFlags.noPrefix)
			})
		} else {
			options := _getFlags.GetCloudWatchLogsGetLogEventsParam()
			if rootFlag.verbose {
				fmt.Printf("GetLogEvents: options=%v\n", options)
			}
			err = cwlogsClient.GetLogEvents(&options, func(out *cloudwatchlogs.GetLogEventsOutput) {
				printEvents(out.Events, !_getFlags.noPrefix)
			})
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	}
}

func init() {
	RootCmd.AddCommand(getCmd)

	localFlags := getCmd.Flags()
	_getFlags.Log.Load(localFlags)
	_getFlags.TimeRange.Load(localFlags)
	localFlags.StringVarP(&_getFlags.Pattern, "filter-pattern", "f", "", "The filter pattern to use. If not provided, all the events are matched.")
	localFlags.BoolVarP(&_getFlags.Watch, "watch", "w", false, "Do not stop when end of log is reached, but rather to wait for additional data to be appended to the input.")
	localFlags.BoolVarP(&_getFlags.noPrefix, "no-prefix", "", false, "Do not display the time and stream name in the event at the begin of the line.")

	// There is no guarantee that _getFlags will have a value here
	// I do not know until getRun function
}
