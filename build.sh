#!/usr/bin/env bash

set -e
set -u

GOARCH=amd64 go build -o sub-amd64
GOARCH=386   go build -o sub-i386
