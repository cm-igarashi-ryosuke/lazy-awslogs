package flags

// goのexport: 大文字から始まる名前がexportされる

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CWLogIdentifyFlags struct {
	Group  string
	Streams []*string
}

func (this *CWLogIdentifyFlags) String() string {
	return fmt.Sprintf("CWLogIdentify={group=%s, streams=%s}", this.Group, this.Streams)
}

// pflag.FlagSetからCWLogIdentifyFlagsをロードする
func (this *CWLogIdentifyFlags) Load(pflag *pflag.FlagSet) {
	pflag.StringVarP(&this.Group, "group", "g", viper.GetString("group"), "The name of log group")
	var strings []string = []string{}
	pflag.StringSliceVarP(&strings, "streams", "s", viper.GetStringSlice("stream"), "The names of log stream")
}
