#!/usr/bin/env bash
protoc --proto_path=./proto --go_opt=Mimage.proto=/ --go_out=./proto ./proto/image.proto
