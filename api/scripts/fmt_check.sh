#!/bin/bash

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}) && cd .. && pwd)"

(
    cd $ROOT_DIR
    VIOLATIONS=$(gofmt -l .)

    if [ ! -z "$VIOLATIONS" ]; then
        echo "Violations found in the following files, run gofmt (or make fmt) locally!"
        echo ""
        echo $VIOLATIONS
        exit 1
    fi
)