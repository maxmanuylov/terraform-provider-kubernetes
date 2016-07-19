package kubernetes

import (
    "github.com/hashicorp/go-uuid"
    "github.com/hashicorp/terraform/helper/schema"
)

func createKubernetesCluster(clusterData *schema.ResourceData, _ interface{}) error {
    id, err := uuid.GenerateUUID()
    if err != nil {
        return err
    }

    client, err := GetOrCreateKubeClient(id, clusterData)
    if err != nil {
        return err
    }

    if err = client.WaitForAPIServer(); err != nil {
        return err
    }

    clusterData.Set("cluster", id)
    clusterData.SetId(id)

    return nil
}

func readKubernetesCluster(clusterData *schema.ResourceData, _ interface{}) error {
    _, err := GetOrCreateKubeClient(clusterData.Id(), clusterData)
    return err
}

func updateKubernetesCluster(clusterData *schema.ResourceData, _ interface{}) error {
    client, err := GetOrCreateKubeClient(clusterData.Id(), clusterData)
    if err != nil {
        return err
    }

    if err = client.WaitForAPIServer(); err != nil {
        return err
    }

    return nil
}

func deleteKubernetesCluster(clusterData *schema.ResourceData, _ interface{}) error {
    DeleteKubeClient(clusterData.Id())

    clusterData.SetId("")
    clusterData.Set("cluster", "")

    return nil
}

func kubernetesClusterExists(clusterData *schema.ResourceData, _ interface{}) (bool, error) {
    _, err := GetOrCreateKubeClient(clusterData.Id(), clusterData)
    return err == nil, err
}
