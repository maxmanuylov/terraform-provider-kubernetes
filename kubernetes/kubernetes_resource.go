package kubernetes

import (
    "bytes"
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/go-rest/client"
    "strings"
)

const (
    namespacesCollection = "namespaces"
    defaultNamespace = "default"
)

type CustomMap map[string]interface{}

func toCustomMap(m interface{}) CustomMap {
    return CustomMap(m.(map[string]interface{}))
}

type KubeResourceId struct {
    apiPath    string
    namespace  string
    collection string
    name       string
}

type KubeResource struct {
    KubeResourceId

    labels      CustomMap
    annotations CustomMap
    waitFor     []interface{}
    content     string
}

func GetKubeResourceId(resourceData *schema.ResourceData) *KubeResourceId {
    resourceNamespace := resourceData.Get("namespace").(string)
    if resourceNamespace == "" {
        resourceNamespace = defaultNamespace
    }

    resourceApiPath := strings.Trim(resourceData.Get("api_path").(string), "/")
    if resourceApiPath == "" {
        resourceApiPath = "api/v1"
    }

    return &KubeResourceId{
        apiPath: resourceApiPath,
        namespace: resourceNamespace,
        collection: resourceData.Get("collection").(string),
        name: resourceData.Get("name").(string),
    }
}

func GetKubeResource(resourceData *schema.ResourceData) *KubeResource {
    return &KubeResource{
        KubeResourceId: *GetKubeResourceId(resourceData),

        labels: toCustomMap(resourceData.Get("labels")),
        annotations: toCustomMap(resourceData.Get("annotations")),
        waitFor: resourceData.Get("wait_for").([]interface{}),
        content: resourceData.Get("content").(string),
    }
}

func (resourceId *KubeResourceId) ApiPath() string {
    return resourceId.apiPath
}

func (resourceId *KubeResourceId) Name() string {
    return resourceId.name
}

func (resourceId *KubeResourceId) IsNamespace() bool {
    return resourceId.collection == namespacesCollection
}

func (resourceId *KubeResourceId) CannotBeDeleted() bool {
    return resourceId.IsNamespace() && (resourceId.name == defaultNamespace || resourceId.name == "kube-system")
}

func (resourceId *KubeResourceId) GetCollection(restClient *rest_client.Client) rest_client.Collection {
    collection := restClient.Collection(fmt.Sprintf("%s/%s", resourceId.apiPath, namespacesCollection))
    if !resourceId.IsNamespace() {
        collection = collection.SubCollection(resourceId.namespace, resourceId.collection)
    }
    return collection
}

func (resourceId *KubeResourceId) Describe() string {
    return fmt.Sprintf("%s/%s/%s/%s", resourceId.apiPath, resourceId.namespace, resourceId.collection, resourceId.name)
}

func (resource *KubeResource) GetWaitFor() []string {
    waitFor := make([]string, len(resource.waitFor))
    for i, wf := range resource.waitFor {
        waitFor[i] = wf.(string)
    }
    return waitFor
}

func (resource *KubeResource) PrepareContent() []byte {
    var buf bytes.Buffer

    buf.WriteString("metadata:\n")
    buf.WriteString("  name: \"")
    buf.WriteString(resource.name)

    if !resource.IsNamespace() {
        buf.WriteString("\"\n  namespace: \"")
        buf.WriteString(resource.namespace)
    }

    buf.WriteString("\"\n")

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
