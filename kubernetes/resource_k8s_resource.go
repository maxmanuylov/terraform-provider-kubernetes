package kubernetes

import (
    "github.com/hashicorp/go-uuid"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/client"
)

func createKubernetesResource(resourceData *schema.ResourceData, _ interface{}) error {
    kubeClient, err := kubernetes_client.GetKubeClient(resourceData)
    if err != nil {
        return err
    }

    id, err := uuid.GenerateUUID()
    if err != nil {
        return err
    }

    if err := kubeClient.Create(kubernetes_client.GetKubeResource(resourceData)); err != nil {
        return err
    }

    resourceData.SetId(id)

    return nil
}

func readKubernetesResource(resourceData *schema.ResourceData, meta interface{}) error {
    exists, err := kubernetesResourceExists(resourceData, meta)
    if err != nil {
        return err
    }

    if !exists {
        resourceData.SetId("")
    }

    return nil
}

func updateKubernetesResource(resourceData *schema.ResourceData, _ interface{}) error {
    kubeClient, err := kubernetes_client.GetKubeClient(resourceData)
    if err != nil {
        return err
    }

    if err := kubeClient.Update(kubernetes_client.GetKubeResource(resourceData)); err != nil {
        return err
    }

    return nil
}

func deleteKubernetesResource(resourceData *schema.ResourceData, _ interface{}) error {
    kubeClient, err := kubernetes_client.GetKubeClient(resourceData)
    if err != nil {
        return err
    }

    if err := kubeClient.Delete(kubernetes_client.GetKubeResourceId(resourceData)); err != nil {
        return err
    }

    resourceData.SetId("")

    return nil
}

func kubernetesResourceExists(resourceData *schema.ResourceData, _ interface{}) (bool, error) {
    kubeClient, err := kubernetes_client.GetKubeClient(resourceData)
    if err != nil {
        return false, err
    }

    return kubeClient.Exists(kubernetes_client.GetKubeResourceId(resourceData))
}
