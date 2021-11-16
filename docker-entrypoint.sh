#!/bin/bash
set -e

if [ "$1" = 'run' ]; then
    sleep 2
    ls -al
    ./accumulate run -n 0
else
    exec "$@"
fi

