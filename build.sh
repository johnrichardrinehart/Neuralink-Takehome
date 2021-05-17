#!/usr/bin/env bash
# build the client and server

function host {
	export GOLOC=/usr/local/go
	export PATH=$GOLOC/bin:$PATH
	go build -o /etc/nl-client ./client
	go build -o /etc/nl-server ./server
}

function docker {
	if [ "$EUID" -eq 0 ]
	then echo "error: \"./build.sh docker\" needs to be run as non-root (without sudo)"
		exit
	fi
	podman build  -t "nl-client:latest" -f Dockerfile.client
	podman build  -t "nl-server:latest" -f Dockerfile.server
}

function main {
	if [[ $1 == "" ]] || [[ $1 == "host" ]] ; then
		host
	elif [[ $1 == "docker" ]]; then
		docker
	fi
}

main $1
