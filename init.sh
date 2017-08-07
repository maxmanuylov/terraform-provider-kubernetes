#!/bin/bash

TERRAFORM_VERSION="v0.10.0"

rm -rf vendor

govendor init

govendor fetch github.com/maxmanuylov/go-rest/client
govendor fetch github.com/maxmanuylov/utils/http/transport/tls
govendor fetch github.com/hashicorp/go-plugin@f72692aebca2008343a9deb06ddb4b17f7051c15
govendor fetch github.com/hashicorp/terraform@=$TERRAFORM_VERSION
govendor fetch github.com/maxmanuylov/utils/intellij-hcl/terraform/provider-schema-generator@=v2
