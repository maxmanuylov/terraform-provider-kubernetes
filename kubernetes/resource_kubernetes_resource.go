package kubernetes

import (
    "github.com/hashicorp/go-uuid"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/client"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/cluster"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/model"
)

func createKubernetesResource(resourceData *schema.ResourceData, meta interface{}) error {
    kubeClient, err := loadClient(resourceData, meta)
    if err != nil {
        return err
    }

    id, err := uuid.GenerateUUID()
    if err != nil {
        return err
    }

    kubeResource, err := kubernetes_model.ParseResource(resourceData)
    if err != nil {
        return err
    }

    if err := kubeClient.Create(kubeResource); err != nil {
        return err
    }

    resourceData.Set("path", kubeResource.Path())
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

func updateKubernetesResource(resourceData *schema.ResourceData, meta interface{}) error {
    kubeClient, err := loadClient(resourceData, meta)
    if err != nil {
        return err
    }

    kubeResource, err := kubernetes_model.ParseResource(resourceData)
    if err != nil {
        return err
    }

    if path := resourceData.Get("path").(string); path == "" {
        if err := kubeClient.Create(kubeResource); err != nil {
            return err
        }

        resourceData.Set("path", kubeResource.Path())
    } else if newPath := kubeResource.Path(); newPath != path {
        if err := kubeClient.Delete(kubernetes_model.ParsePath(path)); err != nil {
            return err
        }

        if err := kubeClient.Create(kubeResource); err != nil {
            return err
        }

        resourceData.Set("path", newPath)
    } else {
        return kubeClient.Update(kubeResource)
    }

    return nil
}

func deleteKubernetesResource(resourceData *schema.ResourceData, meta interface{}) error {
    kubeClient, err := loadClient(resourceData, meta)
    if err != nil {
        return err
    }

    if path := resourceData.Get("path").(string); path != "" {
        if err := kubeClient.Delete(kubernetes_model.ParsePath(path)); err != nil {
            return err
        }
    }

    resourceData.SetId("")

    return nil
}

func kubernetesResourceExists(resourceData *schema.ResourceData, meta interface{}) (bool, error) {
    kubeClient, err := loadClient(resourceData, meta)
    if err != nil {
        return false, err
    }

    if path := resourceData.Get("path").(string); path != "" {
        return kubeClient.Exists(kubernetes_model.ParsePath(path))
    }

    return false, nil
}

func loadClient(resourceData *schema.ResourceData, meta interface{}) (*kubernetes_client.KubeClient, error) {
    cluster, err := kubernetes_cluster.Load(resourceData, meta)
    if err != nil {
        return nil, err
    }
    
    return kubernetes_client.New(cluster)
}
