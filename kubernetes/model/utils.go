package kubernetes_model

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
    "gopkg.in/yaml.v2"
    "strings"
)

type k8sEntity struct {
    ApiVersion string
    Kind       string
    Metadata struct {
        Name      string
        Namespace string
    }
}

func (e *k8sEntity) GetApiPath() string {
    if e.ApiVersion == "" || e.ApiVersion == "v1" {
        return DefaultApiPath
    }
    return fmt.Sprintf("apis/%s", strings.Trim(e.ApiVersion, "/"))
}

func (e *k8sEntity) GetNamespace() string {
    if e.Metadata.Namespace == "" {
        return DefaultNamespace
    }
    return e.Metadata.Namespace
}

func ParseResource(resourceData *schema.ResourceData) (*KubeResource, error) {
    contents := []byte(resourceData.Get("contents").(string))
    encoding := resourceData.Get("encoding").(string)

    entity := &k8sEntity{}
    if encoding == EncodingJson {
        if err := json.Unmarshal(contents, entity); err != nil {
            return nil, err
        }
    } else {
        if err := yaml.Unmarshal(contents, entity); err != nil {
            return nil, err
        }
    }

    if entity.Kind == "" {
        return nil, errors.New("Kind is not specified")
    }

    if entity.Metadata.Name == "" {
        return nil, errors.New("Name is not specified")
    }

    return &KubeResource{
        KubeResourcePath: &KubeResourcePath{
            ApiPath:    entity.GetApiPath(),
            Namespace:  entity.GetNamespace(),
            Collection: fmt.Sprintf("%ss", strings.ToLower(entity.Kind)),
            Name:       entity.Metadata.Name,
        },
        Contents: contents,
        Encoding: encoding,
    }, nil
}

func ParsePath(path string) *KubeResourcePath {
    resourceName, collectionPath := splitOne(path)
    collectionName, restPath := splitOne(collectionPath)

    if collectionName == namespacesCollection {
        return &KubeResourcePath{
            ApiPath:    restPath,
            Namespace:  DefaultNamespace,
            Collection: namespacesCollection,
            Name:       resourceName,
        }
    }

    namespaceName, namespaceCollectionPath := splitOne(restPath)
    _, apiPath := splitOne(namespaceCollectionPath)

    return &KubeResourcePath{
        ApiPath:    apiPath,
        Namespace:  namespaceName,
        Collection: collectionName,
        Name:       resourceName,
    }
}

func splitOne(path string) (string, string) {
    if slashPos := strings.LastIndex(path, "/"); slashPos == -1 {
        return path, ""
    } else {
        return path[slashPos + 1:], path[:slashPos]
    }
}
