#!/bin/bash

TERRAFORM_DIR="<Insert your value>"
LOCAL_OS="<Insert your value>"

cp -f bin/k8s.json "$HOME/.terraform.d/schemas/k8s.json"
cp -f "bin/$LOCAL_OS/terraform-provider-k8s" "$TERRAFORM_DIR/terraform-provider-k8s"
