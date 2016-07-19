#!/bin/bash

TERRAFORM_VERSION="v0.6.16"

rm -rf vendor

govendor init

govendor fetch github.com/hashicorp/terraform@=$TERRAFORM_VERSION
govendor fetch github.com/maxmanuylov/utils/intellij-hcl/terraform/provider-schema-generator@=terraform-$TERRAFORM_VERSION
