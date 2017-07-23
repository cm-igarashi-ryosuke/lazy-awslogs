# lazy-awslogs

[![Join the chat at https://gitter.im/lazyawslogs/Lobby](https://badges.gitter.im/lazyawslogs/Lobby.svg)](https://gitter.im/lazyawslogs/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

"lazy-awslogs" is a tool for Amazon CloudWatch Logs.

This command can query groups and streams, and filter events.
In addition, remember frequently used settings and help to reduce your type amount.

Inspired by [jorgebastida/awslogs: AWS CloudWatch logs for Humans™](https://github.com/jorgebastida/awslogs).

## Demo

![demo](https://github.com/cm-igarashi-ryosuke/lazy-awslogs/blob/master/demo.gif)

## Features

- Show logs, groups, streams.
- Shell completion.
- Configuration.
- Human-friendly time filtering.

## Installing

```
$ go get github.com/cm-igarashi-ryosuke/lazy-awslogs
```

## Shell completion

Shell completion supports your type for subcommands, options, and option's value(e.g. `--group` , `--stream` )

### bash

For Mac:

```
$ cp shell-completion/lazy-awslogs.sh /usr/local/etc/bash_completion.d/lazy-awslogs
$ chmod +x /usr/local/etc/bash_completion.d/lazy-awslogs
$ source /usr/local/etc/bash_completion.d/lazy-awslogs
```

## fish

Add following code in your config.fish

```
$ cd $project
$ source shell-completion/lazy-awslogs.fish
```

## AWS Credentials

- Use `--aws-access-key-id` and `--aws-secret-access-key` options.
- Use `--profile` option.
- Recommending, use `config` command. For details, write below.

## Configuration

`config` command records your settings and define a default setting(Like a rbenv). Then, it be able to omit some options. 

```
$ lazy-awslogs config add \
--config-name    staging \
--config-profile staging-profile \
--config-region  us-east-1

$ lazy-awslogs config add \
--config-name    production \
--config-profile production-profile \
--config-region  us-east-1

$ lazy-awslogs config list
* staging
  production

$ lazy-awslogs config show
====== [Current Environment] ======
EnvironmentName: staging
Profile:         staging-profile
Region:          us-east-1

$ lazy-awslogs config use production

$ lazy-awslogs config list
  staging
* production

$ lazy-awslogs config show
====== [Current Environment] ======
EnvironmentName: production
Profile:         production-profile
Region:          us-east-1
```

Note that `config` command make a configuration file at `~/.lazy-awslogs.yaml` by default.

## Local cache

`reload` command cache all groups and streams. Then, it be able to use `Shell completion` for `--group` and `--stream` options.

Note that depending on the amount, it may take some time.

```
$ lazy-awslogs reload
```

## Usage

For details of each command, please refer to the `--help` .

```
Usage:
  lazy-awslogs [command]

Available Commands:
  config      Manage config file
  get         Lists log events from the specified log stream.
  groups      List groups
  help        Help about any command
  reload      Get groups and streams and update config file
  streams     List streams
  version     Display version

Flags:
      --aws-access-key-id string       AWS access key
      --aws-secret-access-key string   AWS secret key
      --config string                  Config file (default is $HOME/.lazy-awslogs.yaml)
  -h, --help                           help for lazy-awslogs
  -p, --profile string                 AWS credentials profile
  -r, --region string                  AWS region
  -v, --verbose                        Enable verbose messages
```

## TODO

- Test!
- Inplement interactive initialize command for config.
- Implement more option for get command.
- More support credential type.
- And more.

## Dependency Library

- [AWS SDK for Go](https://github.com/aws/aws-sdk-go)
- [cobra](https://github.com/spf13/cobra)
- [YAML support for the Go language](https://github.com/go-yaml/yaml)
- [kindlytime-go](https://github.com/yamaya/kindlytime-go)

## Author
- [cm\-igarashi\-ryosuke](https://github.com/cm-igarashi-ryosuke)

## Contributors
- [yamaya \(Masayuki Yamaya\)](https://github.com/yamaya)
