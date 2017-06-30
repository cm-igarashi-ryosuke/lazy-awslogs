_lazy-awslogs()
{
    local cmd cur prev
    cmd="${COMP_WORDS[1]}"
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    case "${cmd}" in 
        config)
            if [ ${COMP_CWORD} -le 2 ]; then
                local opts="add list remove show set use"
                COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
                return 0
            fi
            return 0
            ;; 
        get)
            if [ "${prev}" = "-g" ] || [ "${prev}" = "--group" ]; then
                COMPREPLY=($(compgen -W "$(get_groups_cache)" ${cur}))
                return 0
            fi
            if [ "${prev}" = "-s" ] || [ "${prev}" = "--stream" ]; then
                COMPREPLY=($(compgen -W "$(get_streams_cache $(current_group_option_value))" -- ${cur}))
                return 0
            fi
            local opts="--filter-pattern --group --help --start-time --stream"
            COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
            return 0
            ;; 
        groups)
            local opts="--cache --help"
            COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
            return 0
            ;;
        reload)
            if [ "${prev}" = "-g" ] || [ "${prev}" = "--group" ]; then
                COMPREPLY=($(compgen -W "$(get_groups_cache)" -- ${cur}))
                return 0
            fi
            local opts="--group --help"
            COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
            return 0
            ;; 
        streams)
            if [ "${prev}" = "-g" ] || [ "${prev}" = "--group" ]; then
                COMPREPLY=($(compgen -W "$(get_groups_cache)" -- ${cur}))
                return 0
            fi
            local opts="--cache --group --help"
            COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
            return 0
            ;;
    esac

    if [ ${COMP_CWORD} -le 2 ]; then
        local opts="config get groups help reload streams"
        COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
        return 0
    fi
}

get_groups_cache()
{
    echo $(lazy-awslogs groups --cache | tr '\n' ' ')
}

get_streams_cache()
{
    if [ $# -eq 0 ]; then
        return
    fi
    echo $(lazy-awslogs streams --cache --group $1 | tr '\n' ' ')
}

current_group_option_value()
{
    local i word group
    i=0
    for word in ${COMP_WORDS[@]}; do
        if [ "${word}" = "-g" ] || [ "${word}" = "--group" ]; then
            group=${COMP_WORDS[${i}+1]}
        fi
        i=${i}+1
    done
    echo ${group}
}

complete -F _lazy-awslogs lazy-awslogs