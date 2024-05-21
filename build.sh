#!/bin/sh
go fmt . ./cmd
go build -o $HOME/.local/bin/docker-cli.exe
