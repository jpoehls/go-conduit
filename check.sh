#!/bin/sh

golint ./
go vet ./
go test