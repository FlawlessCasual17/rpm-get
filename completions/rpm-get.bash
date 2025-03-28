#!/usr/bin/env bash

function _rpm-get() {
    if [ "${COMP_CWORD}" = 1 ]; then
        COMPREPLY=($(compgen -W "update upgrade show install reinstall remove purge search cache clean list pretty_list prettylist csv_list csvlist csv fix-installed help version" "${COMP_WORDS[1]}"))
    elif [ "${COMP_CWORD}" -ge 2 ]; then
        local command="${COMP_WORDS[1]}"

        if  [ "${command}" = update ] &&  [ "${COMP_CWORD}" -le 3 ]; then
            COMPREPLY=($(compgen -W "--repos-only --quiet" "\\${COMP_WORDS[${COMP_CWORD}]}"))
        elif [ "${command}" = show ]; then
            COMPREPLY=($(compgen -W "$(rpm-get list --include-unsupported --raw | tr "\n" " ")" "${COMP_WORDS[${COMP_CWORD}]}"))
        elif [ "${COMP_CWORD}" = 2 ] && [ "${command}" = search ]; then
            COMPREPLY=($(compgen -W "--include-unsupported $(rpm-get list --raw | tr "\n" " ")" "\\${COMP_WORDS[${COMP_CWORD}]}"))
        elif [ "${COMP_CWORD}" = 3 ] && [ "${command}" = search ] && [ "${COMP_WORDS[2]}" = --include-unsupported ]; then
            COMPREPLY=($(compgen -W "$(rpm-get list --include-unsupported --raw | tr "\n" " ")" "${COMP_WORDS[${COMP_CWORD}]}"))
        elif [ "${command}" = install ]; then
            COMPREPLY=($(compgen -W "$(rpm-get list --not-installed | tr "\n" " ")" "${COMP_WORDS[${COMP_CWORD}]}"))
        elif [ "${command}" = reinstall ]; then
            COMPREPLY=($(compgen -W "$(rpm-get list --installed | tr "\n" " ")" "${COMP_WORDS[${COMP_CWORD}]}"))
        elif [[ " remove purge " =~ " ${command} " ]]; then
            if [ "${COMP_CWORD}" = 2 ]; then
                COMPREPLY=($(compgen -W "--remove-repo $(rpm-get list --installed | tr "\n" " ")" "\\${COMP_WORDS[2]}"))
            else
                COMPREPLY=($(compgen -W "$(rpm-get list --installed | tr "\n" " ")" "${COMP_WORDS[${COMP_CWORD}]}"))
            fi
        elif [ "${command}" = list ]; then
            COMPREPLY=($(compgen -W "--include-unsupported --raw --installed --not-installed" "\\${COMP_WORDS[${COMP_CWORD}]}"))
        elif [ "${COMP_CWORD}" = 2 ] && [[ " pretty_list prettylist csv_list csvlist csv " =~ " ${command} " ]]; then
            COMPREPLY=($(compgen -W "$(find "/etc/rpm-get" -maxdepth 1 \( -name *.repo ! -name 00-builtin.repo ! -name 99-local.repo -type f \) -o \( -name 99-local.d -type d \) -printf "%f\n" 2> /dev/null | sed "s/.repo$//; s/.d$//" | tr "\n" " ") 00-builtin 01-main" "${COMP_WORDS[2]}"))
        elif [ "${COMP_CWORD}" = 2 ] && [ "${command}" = fix-installed ]; then
            COMPREPLY=($(compgen -W "--old-apps" "\\${COMP_WORDS[2]}"))
        fi
    fi
}

complete -F _rpm-get rpm-get
