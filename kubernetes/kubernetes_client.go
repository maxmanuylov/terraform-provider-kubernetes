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
var ErrConflict = clientError("409 Conflict")

type KubeClient struct {
    apiUrl     string
    httpClient *http.Client
}

func (client *KubeClient) WaitForAPIServer() error {
    response, err := client.retryLong("GET", "", nil)

    if err != nil {
        dumpErrorToFile("connect to Kubernetes API server", nil, response, err)
    }

    return err
}

func (client *KubeClient) Create(resource *KubeResource) error {
    content := resource.PrepareContent()

    response, err := client.retryLong("POST", resource.GetCollectionPath(), content)

    if err == ErrConflict { // resource already exists
        return client.doUpdate(resource, content)
    }

    if err != nil {
        dumpErrorToFile(fmt.Sprintf("create %q", resource.GetResourcePath()), content, response, err)
    }

    return err
}

func (client *KubeClient) Update(resource *KubeResource) error {
    return client.doUpdate(resource, resource.PrepareContent())
}

func (client *KubeClient) doUpdate(resource *KubeResource, content []byte) error {
    path := resource.GetResourcePath()

    response, err := client.retryLong("PUT", path, content)

    if err != nil {
        dumpErrorToFile(fmt.Sprintf("update %q", path), content, response, err)
    }

    return err
}

func (client *KubeClient) Exists(resourceId *KubeResourceId) (bool, error) {
    response, err := client.retryShort("GET", resourceId.GetResourcePath(), nil)

    if err == ErrNotFound {
        return false, nil
    }

    return err == nil && response.StatusCode / 100 == 2, nil
}

func (client *KubeClient) Delete(resourceId *KubeResourceId) error {
    if resourceId.CannotBeDeleted() {
        return nil
    }

    path := resourceId.GetResourcePath()

    response, err := client.retryShort("DELETE", path, nil)

    if err == ErrNotFound {
        return nil
    }

    if err != nil {
        dumpErrorToFile(fmt.Sprintf("delete %q", path), nil, response, err)
    }

    return err
}

func (client *KubeClient) retryLong(method, path string, content []byte) (*http.Response, error) {
    return client.retry(200, method, path, content) // 10 minutes
}

func (client *KubeClient) retryShort(method, path string, content []byte) (*http.Response, error) {
    return client.retry(3, method, path, content) // 3 times
}

func (client *KubeClient) retry(n int, method, path string, content []byte) (response *http.Response, err error) {
    var done bool
    for i := 0; i < n; i++ {
        if i != 0 {
            time.Sleep(3 * time.Second)
        }

        if done, response, err = client.try(method, path, content); done {
            return
        }

        if err != nil {
            dumpError(os.Stderr, fmt.Sprintf("%s \"%s\"", method, path), content, response, err)
        }
    }
    return
}

func (client *KubeClient) try(method, path string, content []byte) (bool, *http.Response, error) {
    var contentReader io.Reader
    if content != nil {
        contentReader = bytes.NewReader(content)
    }

    response, err := client.do(method, path, contentReader)
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
        if response.StatusCode == 409 {
            return true, response, ErrConflict
        }
        if response.StatusCode == 403 { // Illegal Kubernetes state, need to retry
            return false, response, clientError(response.Status)
        }
        return true, response, clientError(response.Status)
    }

    return false, response, fmt.Errorf("Server error: %s", response.Status)
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

func clientError(status string) error {
    return fmt.Errorf("Client error: %s", status)
}

func dumpErrorToFile(action string, content []byte, response *http.Response, err error) {
    file, err2 := os.Create("kubernetes-error.log")
    if err2 != nil {
        return
    }

    defer file.Close()
    defer file.Sync()

    fmt.Fprintf(file, "%s\n\n", time.Now().String())

    dumpError(file, action, content, response, err)
}

func dumpError(writer io.Writer, action string, content []byte, response *http.Response, err error) {
    fmt.Fprintf(writer, "Failed to %s: %v\n\n", action, err)

    if content != nil {
        fmt.Fprintln(writer, "=============")
        fmt.Fprintln(writer, "== Content ==")
        fmt.Fprintln(writer, "=============")
        fmt.Fprintf(writer, "\n%s\n\n", content)
    }

    if response != nil {
        fmt.Fprintln(writer, "==============")
        fmt.Fprintln(writer, "== Response ==")
        fmt.Fprintln(writer, "==============")
        fmt.Fprintf(writer, "\n%s\n", response.Status)

        if response.Header != nil {
            for key, values := range response.Header {
                for _, value := range values {
                    fmt.Fprintf(writer, "%s: %s\n", key, value)
                }
            }
        }

        fmt.Fprintln(writer)
        io.Copy(writer, response.Body)
        fmt.Fprintln(writer)
    }
}