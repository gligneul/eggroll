#!/bin/bash

usage() {
    echo "usage: scripts/update-examples.sh [-a] [-r rev]

    -a      Add the changes to git
    -r rev  Set eggroll to a specif git revision
    "
}

gitadd=false
gitrev=HEAD

while [[ $# -gt 0 ]]; do
    case $1 in
        -a)
            gitadd=true
            shift
            ;;
        -r)
            shift
            gitrev=$1
            shift
            ;;
        -h)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option $1"
            exit 1
            ;;
    esac
done

rev=$(git rev-parse $gitrev)
if [[ $? -ne 0 ]]; then
    exit 1
fi

for example in ./examples/*; do
    (
        cd $example
        go get github.com/gligneul/eggroll@$rev
        go mod tidy
        if [[ "$gitadd" == "true" ]]; then
            git add go.mod go.sum
        fi
    )
done
