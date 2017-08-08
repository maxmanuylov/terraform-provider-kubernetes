package kubernetes

import (
    "github.com/hashicorp/go-uuid"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/client"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/cluster"
)

func createKubernetesCluster(clusterData *schema.ResourceData, _ interface{}) error {
    id, err := uuid.GenerateUUID()
    if err != nil {
        return err
    }

    cluster := kubernetes_cluster.New(clusterData)

    client, err := kubernetes_client.New(cluster)
    if err != nil {
        return err
    }

    if err = client.WaitForAPIServer(); err != nil {
        return err
    }

    encodedCluster, err := cluster.Encode()
    if err != nil {
        return err
    }

    clusterData.Set("cluster", encodedCluster)
    clusterData.SetId(id)

    return nil
}
