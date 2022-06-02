#!/bin/bash

# delete old ones
for file in ./bin/sasha*; do rm -f $file; done && echo "deleted old versions"
# get version number
VERSION=$(grep -Eo '[0-9]+\.[0-9]+\.[0-9]+' ./version/version.go)

# compile new ones
for release in "darwin amd64" "linux amd64"
do
    set -- $release
    env GOOS="$1" GOARCH="$2" go build -o "bin/sasha_${VERSION}_${1}_${2}" cptdom/sasha
done && echo "New binaries successfully compiled."
