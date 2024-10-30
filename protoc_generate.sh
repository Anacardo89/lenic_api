#!/bin/sh

rm -rf pb
protoc --go_out=. --go-grpc_out=. lenic.proto