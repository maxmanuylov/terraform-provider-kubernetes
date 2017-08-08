package kubernetes_cluster

import (
    "github.com/hashicorp/terraform/helper/schema"
    "encoding/json"
)

type Cluster struct {
    ApiServer  string
    CaCert     string
    ClientCert string
    ClientKey  string
}

func New(clusterData *schema.ResourceData) *Cluster {
    return &Cluster{
        ApiServer:  clusterData.Get("api_server").(string),
        CaCert:     clusterData.Get("ca_cert").(string),
        ClientCert: clusterData.Get("client_cert").(string),
        ClientKey:  clusterData.Get("client_key").(string),
    }
}

func Load(resourceData *schema.ResourceData, meta interface{}) (*Cluster, error) {
    if encodedCluster := resourceData.Get("cluster").(string); encodedCluster != "" {
        return Decode(encodedCluster)
    }
    return meta.(*Cluster), nil
}

func (c *Cluster) Encode() (string, error) {
    encodedCluster, err := json.Marshal(c)
    return string(encodedCluster), err
}

func Decode(encodedCluster string) (*Cluster, error) {
    cluster := &Cluster{}
    return cluster, json.Unmarshal([]byte(encodedCluster), cluster)
}
