#!/usr/bin/env bash
# build the client and server

function host {
	export GOLOC=/usr/local/go
	export PATH=$GOLOC/bin:$PATH
	go build -o /tmp/nl-client ./client
	go build -o /tmp/nl-server ./server
	chown $(logname): /tmp/nl-{client,server}
	ln -s /tmp/nl-client ./nl-client
	ln -s /tmp/nl-server ./nl-server
	chown -h $(logname): ./nl-client ./nl-server
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
	elif [[ $1 == "nix-legacy" ]]; then
		source $HOME/.nix-profile/etc/profile.d/nix.sh
		nix-build
	elif [[ $1 == "nix-flakes" ]]; then
		source $HOME/.nix-profile/etc/profile.d/nix.sh
		nix build
	fi
}

main $1
