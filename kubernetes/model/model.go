package kubernetes_model

import (
    "fmt"
)

const (
    DefaultApiPath       = "api/v1"
    DefaultNamespace     = "default"
    EncodingJson         = "json"
    EncodingYaml         = "yaml"
    namespacesCollection = "namespaces"
)

type KubeResourcePath struct {
    ApiPath    string
    Namespace  string
    Collection string
    Name       string
}

type KubeResource struct {
    *KubeResourcePath

    Contents []byte
    Encoding string
}

func (resourcePath *KubeResourcePath) IsGlobal() bool {
    return resourcePath.Namespace == ""
}

func (resourcePath *KubeResourcePath) IsNamespace() bool {
    return resourcePath.Collection == namespacesCollection
}

func (resourcePath *KubeResourcePath) CannotBeDeleted() bool {
    return resourcePath.IsNamespace() && (resourcePath.Name == DefaultNamespace || resourcePath.Name == "kube-system" || resourcePath.Name == "kube-public")
}

func (resourcePath *KubeResourcePath) CollectionPath() string {
    if resourcePath.IsGlobal() {
        return fmt.Sprintf("%s/%s", resourcePath.ApiPath, resourcePath.Collection)
    }
    return fmt.Sprintf("%s/%s/%s/%s", resourcePath.ApiPath, namespacesCollection, resourcePath.Namespace, resourcePath.Collection)
}

func (resourcePath *KubeResourcePath) Path() string {
    return fmt.Sprintf("%s/%s", resourcePath.CollectionPath(), resourcePath.Name)
}
