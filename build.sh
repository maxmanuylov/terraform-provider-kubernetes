#!/bin/bash

VERSION="v1.3.2"

govendor sync

rm -rf bin

export GO15VENDOREXPERIMENT=1
export GOARCH=amd64

GOOS=darwin  go build -o bin/macos/terraform-provider-k8s
GOOS=linux   go build -o bin/linux/terraform-provider-k8s
GOOS=windows go build -o bin/windows/terraform-provider-k8s.exe

tar czf bin/terraform-provider-k8s-$VERSION-macos.tar.gz --directory=bin/macos terraform-provider-k8s
tar czf bin/terraform-provider-k8s-$VERSION-linux.tar.gz --directory=bin/linux terraform-provider-k8s
zip     bin/terraform-provider-k8s-$VERSION-windows.zip -j bin/windows/terraform-provider-k8s.exe

go run schema/schema.go k8s bin
