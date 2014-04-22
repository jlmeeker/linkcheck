#!/bin/bash

# Setup environment
export GOPATH=`pwd`
export GOBIN=$GOPATH/bin

# Download necessary dependencies
# go get code.google.com/p/gcfg

# Format Go sources
go fmt src/linkcheck/linkcheck.go
go fmt src/linkcheck/dns/dns.go
go fmt src/linkcheck/ping/ping.go


# Build
rm -f $GOBIN/*
rm -rf $GOPATH/pkg/*
go install linkcheck
sudo chown root $GOBIN/linkcheck
sudo chmod 4755 $GOBIN/linkcheck

