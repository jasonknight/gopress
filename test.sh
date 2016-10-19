#!/usr/bin/env bash

set -e
echo "" > coverage.txt
go test -race -coverprofile=profile.out -covermode=atomic $d
go tool cover -func=profile.out
go tool cover -html=profile.out
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
# for d in $(go list ./... | grep -v vendor); do
    
# done
