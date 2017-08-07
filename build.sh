#!/bin/bash

VERSION="v1.2.1"

govendor sync

rm -rf bin

export GO15VENDOREXPERIMENT=1
export GOARCH=amd64

GOOS=darwin  go build -o bin/macos/terraform-provider-kubernetes
GOOS=linux   go build -o bin/linux/terraform-provider-kubernetes
GOOS=windows go build -o bin/windows/terraform-provider-kubernetes.exe

tar czf bin/terraform-provider-kubernetes-$VERSION-macos.tar.gz --directory=bin/macos terraform-provider-kubernetes
tar czf bin/terraform-provider-kubernetes-$VERSION-linux.tar.gz --directory=bin/linux terraform-provider-kubernetes
zip     bin/terraform-provider-kubernetes-$VERSION-windows.zip -j bin/windows/terraform-provider-kubernetes.exe

go run schema/schema.go kubernetes bin
