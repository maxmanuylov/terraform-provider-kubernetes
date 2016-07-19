package main

import (
    "github.com/hashicorp/terraform/plugin"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: kubernetes.Provider,
    })
}