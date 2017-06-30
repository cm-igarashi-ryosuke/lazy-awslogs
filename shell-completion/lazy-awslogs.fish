set __fish_awslogs_command_name 'lazy-awslogs'

function __fish_awslogs_parse_subcommand
	set --local subcmds config get groups streams reload version
	set --local pattern (string join '|' $subcmds)
	for arg in $argv[2..-1]
		if string match --quiet --regex -- "^($pattern)\$" $arg
			echo -n $arg
		end
	end
end

function __fish_awslogs_echo_a
	echo -ne (string join -- '\n' $argv)
end

function __fish_awslogs_get_groups
	printf (_ "%s\tGroup\n") (eval "$__fish_awslogs_command_name groups --cache")
end

function __fish_awslogs_get_streams
	for i in (seq (count $argv))[2..-1] # loop to find `--group` option`
		if string match --quiet --regex -- "^(--group|-g)\$" $argv[$i]
			set --local name $argv[(math $i + 1)]
			printf (_ "%s\tStream\n") (eval "$__fish_awslogs_command_name streams --cache --group $name")
		end
	end
end

set help_option "--help\tDisplay help"
set group_option "--group\tLog Group"
set stream_option "--stream\tLog Stream"

function __fish_awslogs_completion_on_root
	__fish_awslogs_echo_a \
		"config\tManage configuration" \
		"get\tList or filter log events" \
		"groups\tList groups" \
		"reload\tUpdate groups or streams" \
		"streams\tList streams of group" \
		"version\tPrint version" \
		"help\tDisplay help" \
		"--region\tAWS region" \
		"--profile\tAWS profile" \
		"--verbose\tPrint extra debugging information" \
		"--aws-access-key-id\tAWS access key" \
		"--aws-secret-access-key\tAWS secret key" \
		$help_option
	return 0
end

function __fish_awslogs_completion_on_get
	__fish_awslogs_echo_a \
		$group_option \
		$stream_option \
		"--filter-pattern\tSearch pattern" \
		"--start-time\tStart time" \
		"--end-time\tEnd time" \
		$help_option
end

function __fish_awslogs_completion_on_reload
	__fish_awslogs_echo_a $group_option $help_option
end

function __fish_awslogs_completion_on_groups
	__fish_awslogs_echo_a "--cache\tUse cache" $help_option
end

function __fish_awslogs_completion_on_streams
	__fish_awslogs_echo_a $group_option (__fish_awslogs_completion_on_groups)
end

function __fish_awslogs_completion_on_config
	set --local subs_desc \
		"add\tAdd an environment" \
		"list\tList environments" \
		"remove\tRemove an environment" \
		"set\tSet a configuration" \
		"show\tShow current configuration" \
		"use\tUse configuration as current"
	set --local subs
	for desc in $subs_desc
		set items (string split -- '\t' $desc)
		if set --query items[1]
			set subs $subs $items[1]
		end
	end
	set --local pattern (string join '|' $subs)
	set --local sub ''
	for arg in $argv[2..-1]
		if string match --quiet --regex -- "^($pattern)\$" $arg
			set sub $arg
			break
		end
	end
	if test -z $sub
		__fish_awslogs_echo_a $subs_desc
	else
		switch $sub
			case add set
				__fish_awslogs_echo_a \
					'--config-group' '--config-name' '--config-profile' '--config-region' '--config-stream' \
					$help_option
			case '*'
				__fish_awslogs_echo_a $help_option
		end
	end
end

function __fish_awslogs_completion
	if not set --query argv[2]
		__fish_awslogs_completion_on_root
		return
	end
	switch $argv[-1]
		case '--group' '-g'
			__fish_awslogs_get_groups
			return
		case '--stream' '-s'
			__fish_awslogs_get_streams $argv
			return
	end
	set --local sub (__fish_awslogs_parse_subcommand $argv)
	if test -z $sub
		__fish_awslogs_completion_on_root
		return
	end
	switch $sub
		case config; __fish_awslogs_completion_on_config $argv
		case get; __fish_awslogs_completion_on_get
		case reload; __fish_awslogs_completion_on_reload
		case groups; __fish_awslogs_completion_on_groups
		case streams; __fish_awslogs_completion_on_streams
	end
end

function __fish_awslogs_completion_unique
	set --local descs (__fish_awslogs_completion $argv)
	if test (count $argv) -eq 1
		__fish_awslogs_echo_a $descs
	else
		set --local uniques
		set --local tab (echo -ne "\t")
		for desc in $descs
			set t (string split -- $tab $desc)
			set --query t[1]
			and not contains -- $t[1] $argv[2..-1]
			and set uniques $uniques $desc
		end
		test (count $uniques) -eq 0
		and return 1
		__fish_awslogs_echo_a $uniques
	end
	return 0
end

function __fish_awslogs_completion_wrapper
	set --local tokens (commandline -ocp)
	__fish_awslogs_completion_unique $tokens
	return $status
end

complete --command "$__fish_awslogs_command_name" --erase
complete --command "$__fish_awslogs_command_name" \
	--arguments '(__fish_awslogs_completion_wrapper)' \
	--condition '__fish_awslogs_completion_wrapper' \
	--no-files

# test code
#set --local args aws-log config list --help
#__fish_awslogs_completion_unique $args
#printf "\nstatus=%d\n" $status
#set --local results (__fish_awslogs_completion_unique $args)
#printf (_ "[%s]\n") $results
