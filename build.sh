#!/usr/bin/env bash

# https://developers.google.com/protocol-buffers/docs/reference/go-generated stipulates that the ;image
# package specifier for --go_opt should be unnecessary. Apparently it isn't in this use case. So, that's
# why it looks ugly.
protoc --proto_path=./proto \
	--go_out=./proto --go_opt=Mimage.proto="/;image" \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mimage.proto=/ \
	./proto/image.proto