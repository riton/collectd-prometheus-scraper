#!/bin/bash

rm -rf vendor

find . -type f \
 -not -path './vendor/*' \
 -name '*.go' \
 -exec sed -i -e 's,gitlab.in2p3.fr/rferrand,github.com/ccin2p3,g' {} \;

sed -i -e 's,gitlab.in2p3.fr/rferrand,github.com/ccin2p3,g' go.mod

go mod edit -replace collectd.org=github.com/ccin2p3/go-collectd@feature/value_meta

go mod tidy
go mod vendor
