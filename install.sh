#!/bin/bash

TERRAFORM_DIR="<Insert your value>"
LOCAL_OS="<Insert your value>"

cp -f bin/kubernetes.json "$HOME/.terraform.d/schemas/kubernetes.json"
cp -f "bin/$LOCAL_OS/terraform-provider-kubernetes" "$TERRAFORM_DIR/terraform-provider-kubernetes"
