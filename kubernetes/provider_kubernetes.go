package kubernetes

import (
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/hashicorp/terraform/terraform"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/cluster"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/model"
    "strings"
)

func Provider() terraform.ResourceProvider {
    return &schema.Provider{
        Schema:        clusterSchema(false),
        ConfigureFunc: configureKubernetesProvider,

        ResourcesMap: map[string]*schema.Resource{
            "k8s_cluster": {
                Schema: clusterSchema(true),
                Create: createKubernetesCluster,
                Read:   readKubernetesCluster,
                Update: updateKubernetesCluster,
                Delete: deleteKubernetesCluster,
            },

            "k8s_resource": {
                Schema: map[string]*schema.Schema{
                    "cluster": {
                        Type:      schema.TypeString,
                        Sensitive: true,
                        Optional:  true,
                    },
                    "contents": {
                        Type:      schema.TypeString,
                        Required:  true,
                        Sensitive: true,
                    },
                    "encoding": {
                        Type:         schema.TypeString,
                        Optional:     true,
                        Default:      kubernetes_model.EncodingYaml,
                        ValidateFunc: validateResourceEncoding,
                    },
                    "path": {
                        Type:     schema.TypeString,
                        Computed: true,
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

func clusterSchema(clusterResource bool) map[string]*schema.Schema {
    clusterSchema := map[string]*schema.Schema{
        "api_server": {
            Type:     schema.TypeString,
            Required: clusterResource,
            Optional: !clusterResource,
        },
        "ca_cert": {
            Type:      schema.TypeString,
            Optional:  true,
            Sensitive: true,
        },
        "client_cert": {
            Type:      schema.TypeString,
            Optional:  true,
            Sensitive: true,
        },
        "client_key": {
            Type:      schema.TypeString,
            Optional:  true,
            Sensitive: true,
        },
    }

    if clusterResource {
        clusterSchema["cluster"] = &schema.Schema{
            Type:      schema.TypeString,
            Sensitive: true,
            Computed:  true,
        }
    }

    return clusterSchema
}

func validateResourceEncoding(v interface{}, _ string) ([]string, []error) {
    if value := strings.ToLower(v.(string)); value != kubernetes_model.EncodingJson && value != kubernetes_model.EncodingYaml {
        return nil, []error{
            fmt.Errorf("Invalid encoding: %v; possible values are \"%s\" and \"%s\"", v, kubernetes_model.EncodingJson, kubernetes_model.EncodingYaml),
        }
    }
    return nil, nil
}

func configureKubernetesProvider(clusterData *schema.ResourceData) (interface{}, error) {
    return kubernetes_cluster.New(clusterData), nil
}
