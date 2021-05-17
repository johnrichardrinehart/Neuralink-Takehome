#!/usr/bin/env bash

function host {
	apt-get update

	# install protoc
	apt-get install -y unzip

	export PROTOC_FILENAME="protoc-3.12.1-linux-x86_64.zip"
	curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v3.12.1/${PROTOC_FILENAME}" -o $PROTOC_FILENAME
	unzip $PROTOC_FILENAME -d /usr/local
	# clean up
	rm $PROTOC_FILENAME

	# install go1.16.4 (apt doesn't support a new-enough version to support Go modules)
	export GOBIN=/usr/local/go/bin
	export PATH=$GOBIN:$PATH
	export GO_FILENAME="go1.16.4.linux-amd64.tar.gz"
	curl -o $GO_FILENAME https://dl.google.com/go/$GO_FILENAME
	tar -xvf go1.16.4.linux-amd64.tar.gz -C /usr/local
	# clean up
	rm $GO_FILENAME

	# install protoc-gen-go* protoc compiler plugins
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	chown -R $(logname): ~/go

	# compile the proto definition into a Go package
	protoc --proto_path=./proto \
		--go_out=./proto --go_opt=Mimage.proto="/;image" \
		--go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mimage.proto="/;image" \
		./proto/image.proto
}

function docker {
	. /etc/os-release
	echo "deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${VERSION_ID}/ /" | tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
	curl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${VERSION_ID}/Release.key | sudo apt-key add -
	apt-get update
	apt-get -y upgrade
	apt-get install -y podman
}

function main {
	if [[ $1 == "" ]] || [[ $1 == "host" ]] ; then
		host
	elif [[ $1 == "docker" ]]; then
		docker
	elif [[ $1 == "nix-legacy" ]]; then
		curl -L https://nixos.org/nix/install | sh
	elif [[ $1 == "nix-flakes" ]]; then
		curl -L https://nixos.org/nix/install | sh
	fi
}
main $1
