#!/bin/bash

rev=$(git rev-parse ${1:-HEAD})

for example in ./examples/*; do
    (
        cd $example
        go get github.com/gligneul/eggroll@$rev
        go mod tidy
        git add go.mod go.sum
    )
done
