#!/usr/bin/env bash
# build the client and server
export GOLOC=/usr/local/go
export PATH=$GOLOC/bin:$PATH
go build -o /etc/nl-client ./client
go build -o /etc/nl-server ./server
