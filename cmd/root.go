// Copyright © 2017 Ryosuke Igarashi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	Aws "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/aws"
	Config "github.com/cm-igarashi-ryosuke/lazy-awslogs/lib/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootFlag struct {
	cfgFile            string
	profile            string
	region             string
	awsAccessKeyId     string
	awsSecretAccessKey string
	verbose            bool
}

// GetCloudWatchLogsClientParams creates a CloudWatchLogsClientParams by root flags.
func (this *RootFlag) GetCloudWatchLogsClientParams() (params Aws.CloudWatchLogsClientParams) {
	params.Profile = this.profile
	params.Region = this.region
	params.AwsAccessKeyId = this.awsAccessKeyId
	params.AwsSecretAccessKey = this.awsSecretAccessKey
	return
}

var rootFlag RootFlag

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lazy-awslogs",
	Short: "Amazon CloudWatch Logs utility command",
	Long: `"lazy-awslogs" is a tool for Amazon CloudWatch Logs.

This command can query groups and streams, and filter events.
In addition, remember frequently used settings and help to reduce your type amount.`,
	// Prepare to run for persistent
	PersistentPreRun: persistentPreRun,
}

func persistentPreRun(ccmd *cobra.Command, args []string) {
	// Merge config with command line options.
	config := Config.Load()
	environment := config.CurrentEnvironment()
	if rootFlag.profile == "" {
		rootFlag.profile = environment.Profile
	}
	if rootFlag.region == "" {
		rootFlag.region = environment.Region
	}
	if rootFlag.verbose {
		fmt.Printf("Merged rootFlag=%#v\n", rootFlag)
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flag
	persistentFlags := RootCmd.PersistentFlags()
	persistentFlags.StringVarP(&rootFlag.profile, "profile", "p", "", "AWS credentials profile")
	persistentFlags.StringVarP(&rootFlag.region, "region", "r", "", "AWS region")
	persistentFlags.StringVarP(&rootFlag.awsAccessKeyId, "aws-access-key-id", "", "", "AWS access key")
	persistentFlags.StringVarP(&rootFlag.awsSecretAccessKey, "aws-secret-access-key", "", "", "AWS secret key")
	persistentFlags.BoolVarP(&rootFlag.verbose, "verbose", "v", false, "Enable verbose messages")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	persistentFlags.StringVar(&rootFlag.cfgFile, "config", "", "Config file (default is $HOME/.lazy-awslogs.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootFlag.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(rootFlag.cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lazy-awslogs" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(Config.FileName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && rootFlag.verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
