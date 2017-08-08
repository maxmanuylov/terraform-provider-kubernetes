package kubernetes_client

import (
    "fmt"
    "github.com/maxmanuylov/go-rest/client"
    "github.com/maxmanuylov/go-rest/error"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/cluster"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/model"
    "net/http"
    "strings"
    "time"
)

var (
    ErrNotFound = rest_error.NewByCode(http.StatusNotFound)
    ErrConflict = rest_error.NewByCode(http.StatusConflict)
)

type KubeClient struct {
    restClient *rest_client.Client
}

func New(cluster *kubernetes_cluster.Cluster) (*KubeClient, error) {
    transport, err := newTransport(cluster)
    if err != nil {
        return nil, err
    }

    return &KubeClient{
        restClient: rest_client.New(strings.TrimSuffix(cluster.ApiServer, "/"), &http.Client{
            Transport: transport,
            Timeout:   10 * time.Second,
        }),
    }, nil
}

func (client *KubeClient) WaitForAPIServer() error {
    action := "connect to Kubernetes API server"

    eh := retryLong(action, nil, func() error {
        _, err := client.restClient.Do("GET", kubernetes_model.DefaultApiPath, rest_client.Json, nil)
        return err
    })

    if eh.error != nil {
        dumpErrorsToFile(action, nil, eh)
    }

    return eh.error
}

func (client *KubeClient) Create(resource *kubernetes_model.KubeResource) error {
    collection := client.restClient.Collection(resource.CollectionPath())
    action := fmt.Sprintf("create %s", resource.Path())

    eh := retryLong(action, resource.Contents, func() error {
        err := createResource(collection, resource.Encoding, resource.Contents)
        if err == http.ErrNoLocation {
            return nil
        }
        return err
    })

    if eh.error == ErrConflict { // resource already exists
        return client.Update(resource)
    }

    if eh.error != nil {
        dumpErrorsToFile(action, resource.Contents, eh)
        return eh.error
    }

    return nil
}

func (client *KubeClient) Update(resource *kubernetes_model.KubeResource) error {
    collection := client.restClient.Collection(resource.CollectionPath())
    action := fmt.Sprintf("update %s", resource.Path())

    eh := retryLong(action, resource.Contents, func() error {
        return updateResource(collection, resource.Name, resource.Encoding, resource.Contents)
    })

    if eh.error != nil {
        dumpErrorsToFile(action, resource.Contents, eh)
    }

    return eh.error
}

func (client *KubeClient) Exists(resourcePath *kubernetes_model.KubeResourcePath) (bool, error) {
    collection := client.restClient.Collection(resourcePath.CollectionPath())
    action := fmt.Sprintf("check existence of %s", resourcePath.Path())

    exists := false
    eh := retryShort(action, nil, func() error {
        var err error
        exists, err = collection.Exists(resourcePath.Name)
        return err
    })

    return exists, eh.error
}

func (client *KubeClient) Delete(resourcePath *kubernetes_model.KubeResourcePath) error {
    if resourcePath.CannotBeDeleted() {
        return nil
    }

    collection := client.restClient.Collection(resourcePath.CollectionPath())
    action := fmt.Sprintf("delete %s", resourcePath.Path())

    eh := retryShort(action, nil, func() error {
        return collection.Delete(resourcePath.Name)
    })

    if eh.error == ErrNotFound {
        return nil
    }

    if eh.error != nil {
        dumpErrorsToFile(action, nil, eh)
    }

    return eh.error
}
