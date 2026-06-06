#!/bin/bash

# If a command was passed to the container run the binary 
if [ "$#" != "0" ]; then
    go run ./cmd/server/main.go "$@"
    exit "$?"
fi

DEBUG_OPT=""

if [ "$AIR_DEBUG" == "true" ]; then
    DEBUG_OPT="-d"
fi

if [ -z "$APP_CMD" ]; then
    echo "Must supply APP_CMD env var"
    exit 1
fi 

air $DEBUG_OPT -build.bin ./bin/panoptes-server -build.cmd "./scripts/build.sh" -build.args_bin "$APP_CMD"