package kubernetes

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
)

var ErrNotFound = clientError("404 Not Found")

type KubeClient struct {
    apiUrl     string
    httpClient *http.Client
}

func (client *KubeClient) WaitForAPIServer() error {
    response, err := retry(func() (*http.Response, error) {
        return client.do("GET", "", nil)
    })

    if err != nil {
        dumpError("connect to Kubernetes API server", nil, response, err)
    }

    return err
}

func (client *KubeClient) Create(resource *KubeResource) error {
    content := resource.PrepareContent(true)
    path := resource.GetCollectionPath()

    response, err := retry(func() (*http.Response, error) {
        return client.do("POST", path, bytes.NewReader(content))
    })

    if err != nil {
        dumpError(fmt.Sprintf("create %q", resource.GetResourcePath()), content, response, err)
    }

    return err
}

func (client *KubeClient) Patch(resource *KubeResource) error {
    content := resource.PrepareContent(false)
    path := resource.GetResourcePath()

    response, err := retry(func() (*http.Response, error) {
        return client.do("PATCH", path, bytes.NewReader(content))
    })

    if err != nil {
        dumpError(fmt.Sprintf("update %q", resource.GetResourcePath()), content, response, err)
    }

    return err
}

func (client *KubeClient) Exists(resourceId *KubeResourceId) (bool, error) {
    path := resourceId.GetResourcePath()

    response, err := retry(func() (*http.Response, error) {
        return client.do("GET", path, nil)
    })

    if err == ErrNotFound {
        return false, nil
    }

    return err == nil && response.StatusCode / 100 == 2, err
}

func (client *KubeClient) Delete(resourceId *KubeResourceId) error {
    path := resourceId.GetResourcePath()

    response, err := retry(func() (*http.Response, error) {
        return client.do("DELETE", path, nil)
    })

    if err == ErrNotFound {
        return nil
    }

    if err != nil {
        dumpError(fmt.Sprintf("delete %q", resourceId.GetResourcePath()), nil, response, err)
    }

    return err
}

func (client *KubeClient) do(method, path string, body io.Reader) (*http.Response, error) {
    url := fmt.Sprintf("%s/%s", strings.TrimSuffix(client.apiUrl, "/"), path)

    request, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    request.Header.Add("Content-Type", "application/yaml")
    request.Header.Add("Accept", "application/yaml")
    request.Header.Add("User-Agent", "curl/7.43.0")

    return client.httpClient.Do(request)
}

func retry(do func() (*http.Response, error)) (response *http.Response, err error) {
    var done bool
    for i := 0; i < 200; i++ { // 10 minutes
        if done, response, err = try(do); done {
            return
        }
        time.Sleep(3 * time.Second)
    }
    return
}

func try(do func() (*http.Response, error)) (bool, *http.Response, error) {
    response, err := do()
    if err != nil {
        return false, nil, err
    }

    responseKind := response.StatusCode / 100

    if responseKind == 2 {
        return true, response, nil
    }

    if responseKind == 4 {
        if response.StatusCode == 404 {
            return true, response, ErrNotFound
        }
        if response.StatusCode == 403 { // Illegal Kubernetes state, need to retry
            return false, response, clientError(response.Status)
        }
        return true, response, clientError(response.Status)
    }

    return false, response, fmt.Errorf("Server error: %s", response.Status)
}

func clientError(status string) error {
    return fmt.Errorf("Client error: %s", status)
}

func dumpError(action string, content []byte, response *http.Response, err error) {
    file, err2 := os.Create("kubernetes-error.log")
    if err2 != nil {
        return
    }

    defer file.Close()

    fmt.Fprintf(file, "%s\n\nFailed to %s: %v\n\n", time.Now().String(), action, err)

    if content != nil {
        fmt.Fprintf(file, "=============\n")
        fmt.Fprintf(file, "== Content ==\n")
        fmt.Fprintf(file, "=============\n")
        fmt.Fprintf(file, "\n%s\n\n", content)
    }

    if response != nil {
        fmt.Fprintf(file, "==============\n")
        fmt.Fprintf(file, "== Response ==\n")
        fmt.Fprintf(file, "==============\n")
        fmt.Fprintf(file, "\n%s\n", response.Status)

        if response.Header != nil {
            for key, values := range response.Header {
                for _, value := range values {
                    fmt.Fprintf(file, "%s: %s\n", key, value)
                }
            }
        }

        fmt.Fprintf(file, "\n")
        io.Copy(file, response.Body)
        fmt.Fprintf(file, "\n")
    }

    if err = file.Sync(); err != nil {
        return
    }
}