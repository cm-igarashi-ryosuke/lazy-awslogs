package flags

// goのexport: 大文字から始まる名前がexportされる

import (
	"fmt"
	Pflag "github.com/spf13/pflag"
	KindlyTime "github.com/yamaya/kindlytime-go"
)

type CWLogTimeRange struct {
	startString string
	endString string
}

func (this *CWLogTimeRange) String() string {
	return fmt.Sprintf("CWLogTimeRange={%s,%s}", this.startString, this.endString)
}

func (this *CWLogTimeRange) StartTimeMilliseconds() int64 {
	time, err := KindlyTime.ParseBaseOnCurrentTime(this.startString)
	if err != nil {
		fmt.Errorf("Unrecognized start-time option -- %s\n", err.Error())
		return 0
	}
	return time.Unix() * 1000
}

func (this *CWLogTimeRange) EndTimeMilliseconds() int64 {
	time, err := KindlyTime.ParseBaseOnCurrentTime(this.endString)
	if err != nil {
		fmt.Errorf("Unrecognized end-time option -- %s\n", err.Error())
		return 0
	}
	return time.Unix() * 1000
}

// pflag.FlagSetからCWLogTimeRangeをロードする
func (this *CWLogTimeRange) Load(pflag *Pflag.FlagSet) {
	pflag.StringVarP(&this.startString, "start-time", "", "30 minutes ago", "The start of the time range, expressed as the number of milliseconds since Jan 1, 1970 00:00:00 UTC. Events with a timestamp earlier than this time are not included.")
	pflag.StringVarP(&this.endString, "end-time", "", "now", "The end of the time range, expressed as the number of milliseconds since Jan 1, 1970 00:00:00 UTC. Events with a timestamp later than this time are not included.")
}
