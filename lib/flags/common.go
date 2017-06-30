package flags

// goのexport: 大文字から始まる名前がexportされる

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CWLogIdentifyFlags struct {
	Group  string
	Stream string
}

func (this *CWLogIdentifyFlags) String() string {
	return fmt.Sprintf("CWLogIdentify={group=%s, stream=%s}", this.Group, this.Stream)
}

// pflag.FlagSetからCWLogIdentifyFlagsをロードする
func (this *CWLogIdentifyFlags) Load(pflag *pflag.FlagSet) {
	pflag.StringVarP(&this.Group, "group", "g", viper.GetString("group"), "The name of log group")
	pflag.StringVarP(&this.Stream, "stream", "s", viper.GetString("stream"), "The name of log stream")
}
