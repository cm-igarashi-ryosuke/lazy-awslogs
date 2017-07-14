package aws

import (
	Aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// CloudWatchLogsClient wraps CloudWatchLogs.
//
// FIXME: If it is possible, fix like `type CloudWatchLogsClient *cloudwatchlogs.CloudWatchLogs`
type CloudWatchLogsClient struct {
	cwLogs *cloudwatchlogs.CloudWatchLogs
}

// CloudWatchLogsClientParams is paramater for initialize CloudWatchLogsClient.
type CloudWatchLogsClientParams struct {
	Region             string
	Profile            string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
}

// NewCloudWatchLogsClient creates a CloudWatchLogsClient object with a CloudWatchLogsClientParams.
func NewCloudWatchLogsClient(params CloudWatchLogsClientParams) CloudWatchLogsClient {
	var sess *session.Session
	if params.Profile != "" && params.Region == "" {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           params.Profile,
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		config := Aws.Config{Region: &params.Region}
		if params.AwsAccessKeyId != "" && params.AwsSecretAccessKey != "" {
			config.Credentials = credentials.NewStaticCredentialsFromCreds(credentials.Value{
				AccessKeyID:     params.AwsAccessKeyId,
				SecretAccessKey: params.AwsSecretAccessKey,
			})
		}
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			// If both Config and Profile are present, Config will be given preference.
			Config:            config,
			Profile:           params.Profile,
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	return CloudWatchLogsClient{cwLogs: cloudwatchlogs.New(sess)}
}

// GetLogGroups gets all LogGroup, then processes by argument f.
func (this *CloudWatchLogsClient) GetLogGroups(f func(*cloudwatchlogs.DescribeLogGroupsOutput)) (err error) {
	err = this.cwLogs.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{
		Limit: Aws.Int64(50),
	}, func(out *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
		f(out)
		return !lastPage
	})
	return
}

// GetLogStreams gets LogStream, then processes by argument f.
func (this *CloudWatchLogsClient) GetLogStreams(group string, f func(*cloudwatchlogs.DescribeLogStreamsOutput)) (err error) {
	err = this.cwLogs.DescribeLogStreamsPages(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: Aws.String(group),
		Limit:        Aws.Int64(50),
	}, func(out *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
		f(out)
		return !lastPage
	})
	return
}

// GetLogEvents
func (this *CloudWatchLogsClient) GetLogEvents(input *cloudwatchlogs.GetLogEventsInput, closure func(*cloudwatchlogs.GetLogEventsOutput)) (err error) {
	var tokenPtr *string
	err = this.cwLogs.GetLogEventsPages(input, func(out *cloudwatchlogs.GetLogEventsOutput, lastPage bool) bool {
		if tokenPtr == nil || *tokenPtr != *out.NextForwardToken {
			closure(out)
			tokenPtr = out.NextForwardToken
			return true
		}
		return false
	})
	return
}

// FilterLogEvents
func (this *CloudWatchLogsClient) FilterLogEvents(input *cloudwatchlogs.FilterLogEventsInput, closure func(*cloudwatchlogs.FilterLogEventsOutput)) (err error) {
	var tokenPtr *string
	err = this.cwLogs.FilterLogEventsPages(input, func(out *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		if tokenPtr == nil || out.NextToken == nil || *tokenPtr != *out.NextToken {
			closure(out)
			if out.NextToken != nil {
				tokenPtr = out.NextToken
			}
			return true
		}
		return false
	})
	return
}
