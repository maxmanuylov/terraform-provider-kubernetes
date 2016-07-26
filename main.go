package main

import (
    "github.com/hashicorp/terraform/plugin"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes"
    //"log"
    //"net/http"
    //_ "net/http/pprof"
)

func main() {
    //go func() {
    //    log.Println(http.ListenAndServe("localhost:9999", nil))
    //}()

    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: kubernetes.Provider,
    })
}