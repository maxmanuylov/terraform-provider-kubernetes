#!/bin/bash

VERSION="v1.1"

rm -rf bin

export GO15VENDOREXPERIMENT=1

GOOS=darwin  GOARCH=amd64 go build -o bin/macos/terraform-provider-kubernetes
GOOS=linux   GOARCH=amd64 go build -o bin/linux/terraform-provider-kubernetes
GOOS=windows GOARCH=amd64 go build -o bin/windows/terraform-provider-kubernetes.exe

tar czf bin/terraform-provider-kubernetes-$VERSION-macos.tar.gz --directory=bin/macos terraform-provider-kubernetes
tar czf bin/terraform-provider-kubernetes-$VERSION-linux.tar.gz --directory=bin/linux terraform-provider-kubernetes
zip     bin/terraform-provider-kubernetes-$VERSION-windows.zip -j bin/windows/terraform-provider-kubernetes.exe

go run schema/schema.go kubernetes bin
