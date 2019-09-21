#!/bin/bash
case "$1" in
    "-build") # -build main input output
        shift
        mkdir /exec
        if test "$1" == "-git"
        then mkdir -p /workspace/source
             git clone "$2" /workspace/source
             shift
             shift
        fi
        /bin/proxy -compile "$1"  <"$2" >/exec/exec.zip
        env HOME=/root /bin/ftl.par \
        --base "$3" \
        --name "$4" \
        --directory /exec \
        --destination /exec
    ;;
    *)
        if test -e /exec/exec.zip
        then exec env OW_AUTOINIT=/exec/exec.zip /bin/proxy
        else exec /bin/proxy "$@"
        fi
    ;;
esac
