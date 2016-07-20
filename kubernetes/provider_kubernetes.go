package kubernetes

import (
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/hashicorp/terraform/terraform"
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
                        ForceNew: true,
                    },
                    "api_version": {
                        Type: schema.TypeString,
                        Optional: true,
                    },
                    "ca_cert": {
                        Type: schema.TypeString,
                        Optional: true,
                    },
                    "client_cert": {
                        Type: schema.TypeString,
                        Optional: true,
                    },
                    "client_key": {
                        Type: schema.TypeString,
                        Optional: true,
                    },
                    "cluster": {
                        Type: schema.TypeString,
                        Computed: true,
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
                        ForceNew: true,
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
                    "content": {
                        Type: schema.TypeString,
                        Required: true,
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

