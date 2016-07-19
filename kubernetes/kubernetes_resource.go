package kubernetes

import (
    "bytes"
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
)

type CustomMap map[string]interface{}

func toCustomMap(m interface{}) CustomMap {
    return CustomMap(m.(map[string]interface{}))
}

type KubeResourceId struct {
    namespace  string
    collection string
    name       string
}

type KubeResource struct {
    KubeResourceId

    labels      CustomMap
    annotations CustomMap
    content     string
}

func GetKubeResourceId(resourceData *schema.ResourceData) *KubeResourceId {
    return &KubeResourceId{
        namespace: resourceData.Get("namespace").(string),
        collection: resourceData.Get("collection").(string),
        name: resourceData.Get("name").(string),
    }
}

func GetKubeResource(resourceData *schema.ResourceData) *KubeResource {
    return &KubeResource{
        KubeResourceId: *GetKubeResourceId(resourceData),

        labels: toCustomMap(resourceData.Get("labels")),
        annotations: toCustomMap(resourceData.Get("annotations")),
        content: resourceData.Get("content").(string),
    }
}

func (resourceId *KubeResourceId) GetCollectionPath() string {
    if resourceId.collection == "namespaces" {
        return "namespaces"
    }
    return fmt.Sprintf("namespaces/%s/%s", resourceId.namespace, resourceId.collection)
}

func (resourceId *KubeResourceId) GetResourcePath() string {
    return fmt.Sprintf("%s/%s", resourceId.GetCollectionPath(), resourceId.name)
}

func (resource *KubeResource) PrepareContent(includeNameData bool) []byte {
    var buf bytes.Buffer

    if includeNameData || len(resource.labels) != 0 || len(resource.annotations) != 0 {
        buf.WriteString("metadata:\n")
    }

    if includeNameData {
        buf.WriteString("  name: \"")
        buf.WriteString(resource.name)
        buf.WriteString("\"\n  namespace: \"")
        buf.WriteString(resource.namespace)
        buf.WriteString("\"\n")
    }

    if len(resource.labels) != 0 {
        writeMap(&buf, "labels", resource.labels)
    }

    if len(resource.annotations) != 0 {
        writeMap(&buf, "annotations", resource.annotations)
    }

    buf.WriteString(resource.content)

    return buf.Bytes()
}

func writeMap(buffer *bytes.Buffer, label string, _map CustomMap) {
    buffer.WriteString("  ")
    buffer.WriteString(label)
    buffer.WriteString(":\n")

    for key, value := range _map {
        buffer.WriteString("    ")
        buffer.WriteString(key)
        buffer.WriteString(": \"")
        buffer.WriteString(fmt.Sprintf("%v", value))
        buffer.WriteString("\"\n")
    }
}
