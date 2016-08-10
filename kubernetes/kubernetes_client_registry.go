package kubernetes

import (
    "crypto/tls"
    "crypto/x509"
    "errors"
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/maxmanuylov/go-rest/client"
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

    apiVersion := clusterData.Get("api_version").(string)
    if apiVersion == "" {
        apiVersion = "v1"
    }

    caCert := clusterData.Get("ca_cert").(string)
    clientCert := clusterData.Get("client_cert").(string)
    clientKey := clusterData.Get("client_key").(string)

    transport := http.DefaultTransport.(*http.Transport)

    tlsConfig, err := newTLSConfig(caCert, clientCert, clientKey)
    if err != nil {
        return nil, err
    }

    if tlsConfig != nil {
        transport.TLSClientConfig = tlsConfig
    }

    return &KubeClient{
        restClient: rest_client.New(
            fmt.Sprintf("%s/api/%s", strings.TrimSuffix(apiServer, "/"), apiVersion),
            &http.Client{
                Transport: transport,
                Timeout: 10 * time.Second,
            },
        ),
    }, nil
}

func newTLSConfig(caCert, clientCert, clientKey string) (*tls.Config, error) {
    if caCert == "" || clientCert == "" || clientKey == "" {
        return nil, nil
    }

    caPool := x509.NewCertPool()
    if !caPool.AppendCertsFromPEM([]byte(caCert)) {
        return nil, errors.New("No CA certificate found")
    }

    certificate, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
    if err != nil {
        return nil, err
    }

    return &tls.Config{
        Certificates: []tls.Certificate{certificate},
        RootCAs: caPool,
    }, nil
}
