#!/usr/bin/env bash
rm -rf ./out
go build -o ./out/server.exe server/server.go
go build -o ./out/client.exe client/client.go
