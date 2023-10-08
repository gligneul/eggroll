#!/bin/bash

usage() {
    echo "usage: scripts/test-examples.sh [-c]

    -c  clean the sunodo build after each test"
}

clean=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -c)
            clean=true
            shift
            ;;
        -h)
            usage
            exit 0
            ;;
        -*)
            echo "Unknown option $1"
            exit 1
            ;;
    esac
done


for example in ./examples/*; do
    (
        cd $example

        if [ "$clean" = "true" -a -d .sunodo ]; then
            rm -r .sunodo
        fi

        go test -v || exit 1

        if [ "$clean" = "true" -a -d .sunodo ]; then
            rm -r .sunodo
        fi
    ) || exit 1
done
