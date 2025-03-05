#!/bin/bash

GOOS=js GOARCH=wasm go build -ldflags '-s -w' -o pkg/assets/static/rogueV3.wasm rogue/v3/wasm/main.go
go build -tags netgo -ldflags '-s -w' -o app cmd/reviso/main.go
