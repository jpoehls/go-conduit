#!/bin/sh

go build && go vet ./ && golint ./ && go test