package kubernetes

import (
    "github.com/hashicorp/go-uuid"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/client"
)

func createKubernetesCluster(clusterData *schema.ResourceData, _ interface{}) error {
    id, err := uuid.GenerateUUID()
    if err != nil {
        return err
    }

    name := clusterData.Get("name").(string)

    client, err := kubernetes_client.GetOrCreateKubeClient(name, clusterData)
    if err != nil {
        return err
    }

    if err = client.WaitForAPIServer(); err != nil {
        return err
    }

    clusterData.Set("cluster", name)
    clusterData.SetId(id)

    return nil
}
