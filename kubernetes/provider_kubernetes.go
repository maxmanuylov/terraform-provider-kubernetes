package kubernetes

import (
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/hashicorp/terraform/terraform"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/client"
)

func Provider() terraform.ResourceProvider {
    return &schema.Provider{
        Schema: map[string]*schema.Schema{},

        ResourcesMap: map[string]*schema.Resource{

            "kubernetes_cluster": {
                Schema: map[string]*schema.Schema{
                    "api_server": {
                        Type: schema.TypeString,
                        Required: true,
                    },
                    "ca_cert": {
                        Type: schema.TypeString,
                        Optional: true,
                        Sensitive: true,
                    },
                    "client_cert": {
                        Type: schema.TypeString,
                        Optional: true,
                        Sensitive: true,
                    },
                    "client_key": {
                        Type: schema.TypeString,
                        Optional: true,
                        Sensitive: true,
                    },
                    "cluster": {
                        Type: schema.TypeString,
                        Computed: true,
                        Sensitive: true,
                    },
                },
                Create: createKubernetesCluster,
                Read:   readKubernetesCluster,
                Update: updateKubernetesCluster,
                Delete: deleteKubernetesCluster,
                Exists: kubernetesClusterExists,
            },

            "kubernetes_resource": {
                Schema: map[string]*schema.Schema{
                    "cluster": {
                        Type: schema.TypeString,
                        Required: true,
                        Sensitive: true,
                    },
                    "api_path": {
                        Type: schema.TypeString,
                        Optional: true,
                        Default: kubernetes_client.DefaultApiPath,
                    },
                    "namespace": {
                        Type: schema.TypeString,
                        Optional: true,
                        ForceNew: true,
                    },
                    "collection": {
                        Type: schema.TypeString,
                        Required: true,
                        ForceNew: true,
                    },
                    "name": {
                        Type: schema.TypeString,
                        Required: true,
                        ForceNew: true,
                    },
                    "labels": {
                        Type: schema.TypeMap,
                        Optional: true,
                    },
                    "annotations": {
                        Type: schema.TypeMap,
                        Optional: true,
                    },
                    "wait_for": {
                        Type: schema.TypeList,
                        Optional: true,
                        Elem: &schema.Schema{
                            Type: schema.TypeString,
                        },
                    },
                    "content": {
                        Type: schema.TypeString,
                        Required: true,
                        Sensitive: true,
                    },
                },
                Create: createKubernetesResource,
                Read:   readKubernetesResource,
                Update: updateKubernetesResource,
                Delete: deleteKubernetesResource,
                Exists: kubernetesResourceExists,
            },

        },

    }
}

