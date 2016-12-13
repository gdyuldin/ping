# Introduction

This project is a simple `ping` analog with configurable ICMP `id` field. It can be used with CirrOS to OpenStack network testing.

## Build

    $ env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v ping.go


## Comressing

    $ go get github.com/pwaller/goupx
    $ goupx --best ping


## Run

    # sudo GOPATH=$GOPATH go run ping.go <args>
