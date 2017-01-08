package kubernetes_client

import (
    "errors"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/go-rest/client"
    "github.com/maxmanuylov/utils/http/transport/tls"
    "net/http"
    "strings"
    "time"
)

var registry = make(map[string]*KubeClient)

var ErrKubeClientNotFound = errors.New("Kubernetes client is not found")

func GetKubeClient(resourceData *schema.ResourceData) (*KubeClient, error) {
    return doGetOrCreateKubeClient(resourceData.Get("cluster").(string), nil)
}

func GetOrCreateKubeClient(clusterKey string, clusterData *schema.ResourceData) (*KubeClient, error) {
    return doGetOrCreateKubeClient(clusterKey, func() (*KubeClient, error) {
        return newKubeClient(clusterData)
    })
}

func DeleteKubeClient(clusterKey string) {
    delete(registry, clusterKey)
}

func doGetOrCreateKubeClient(clusterKey string, factory func() (*KubeClient, error)) (*KubeClient, error) {
    client, ok := registry[clusterKey]
    if ok {
        return client, nil
    }

    if factory != nil {
        client, err := factory()
        if err == nil {
            registry[clusterKey] = client
        }
        return client, err
    }

    return nil, ErrKubeClientNotFound
}

func newKubeClient(clusterData *schema.ResourceData) (*KubeClient, error) {
    apiServer := clusterData.Get("api_server").(string)

    caCert := clusterData.Get("ca_cert").(string)
    clientCert := clusterData.Get("client_cert").(string)
    clientKey := clusterData.Get("client_key").(string)

    transport := http.DefaultTransport
    if caCert != "" && clientCert != "" && clientKey != "" {
        var err error
        transport, err = tls_transport.New([]byte(caCert), []byte(clientCert), []byte(clientKey))
        if err != nil {
            return nil, err
        }
    }

    return &KubeClient{
        restClient: rest_client.New(
            strings.TrimSuffix(apiServer, "/"),
            &http.Client{
                Transport: transport,
                Timeout: 10 * time.Second,
            },
        ),
    }, nil
}
