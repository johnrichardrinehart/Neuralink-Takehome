#!/usr/bin/env bash
protoc --proto_path=./proto \
	--go_out=./proto --go_opt=Mimage.proto=/ \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mimage.proto=/ \
	./proto/image.proto
