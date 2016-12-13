#!/usr/bin/env bash

env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
goupx --best ping
