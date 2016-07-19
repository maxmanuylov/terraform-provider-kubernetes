package main

import (
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes"
    "github.com/maxmanuylov/utils/intellij-hcl/terraform/provider-schema-generator"
)

func main() {
    provider_schema_generator.Generate(kubernetes.Provider().(*schema.Provider))
}
