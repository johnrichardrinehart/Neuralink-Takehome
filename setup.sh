#!/usr/bin/env bash

function host {
	apt-get update
	# install protoc
	apt-get install unzip
	export PROTOC_FILENAME="protoc-3.12.1-linux-x86_64.zip"
	curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v3.12.1/${PROTOC_FILENAME}"
	unzip $PROTOC_FILENAME -d /usr/local
	rm $PROTOC_FILENAME
	# install go1.15.2 (apt doesn't support a new-enough version to support Go modules)
	export GO_FILENAME="go1.15.2.linux-amd64.tar.gz"
	wget https://dl.google.com/go/$GO_FILENAME
	tar -xvf go1.15.2.linux-amd64.tar.gz -C /usr/local
	rm $GO_FILENAME
	export GOBIN=/usr/local/go/bin
	export PATH=$GOBIN:$PATH
	# install protoc-gen-go* protoc compiler plugins
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	# compile the proto definition into a Go package
	protoc --proto_path=./proto \
		--go_out=./proto --go_opt=Mimage.proto="/;image" \
		--go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mimage.proto="/;image" \
		./proto/image.proto
}

function main {
	if [[ $1 == "" ]] || [[ $1 == "host" ]] ; then
		host
	fi
}

main $1
