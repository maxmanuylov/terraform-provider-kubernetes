package kubernetes

import (
    "fmt"
    "github.com/maxmanuylov/go-rest/client"
    "github.com/maxmanuylov/go-rest/error"
    "io"
    "net/http"
    "os"
    "time"
)

var ErrNotFound = rest_error.NewByCode(http.StatusNotFound)
var ErrConflict = rest_error.NewByCode(http.StatusConflict)

type KubeClient struct {
    restClient *rest_client.Client
}

func (client *KubeClient) WaitForAPIServer() error {
    action := "connect to Kubernetes API server"

    err := retryLong(action, nil, func() error {
        _, err := client.restClient.Do("GET", "", rest_client.Json, nil)
        return err
    })

    if err != nil {
        dumpErrorToFile(action, nil, err)
    }

    return err
}

func (client *KubeClient) Create(resource *KubeResource) error {
    collection := resource.GetCollection(client.restClient)

    action := fmt.Sprintf("create %s", resource.Describe())
    content := resource.PrepareContent()

    err := retryLong(action, content, func() error {
        _, err := collection.CreateYaml(content)
        if err == http.ErrNoLocation {
            return nil
        }
        return err
    })

    if err == ErrConflict { // resource already exists
        return client.doUpdate(resource, content)
    }

    if err != nil {
        dumpErrorToFile(action, content, err)
    }

    return err
}

func (client *KubeClient) Update(resource *KubeResource) error {
    return client.doUpdate(resource, resource.PrepareContent())
}

func (client *KubeClient) doUpdate(resource *KubeResource, content []byte) error {
    collection := resource.GetCollection(client.restClient)
    action := fmt.Sprintf("update %s", resource.Describe())

    err := retryLong(action, content, func() error {
        return collection.ReplaceYaml(resource.Name(), content)
    })

    if err != nil {
        dumpErrorToFile(action, content, err)
    }

    return err
}

func (client *KubeClient) Exists(resourceId *KubeResourceId) (bool, error) {
    collection := resourceId.GetCollection(client.restClient)
    action := fmt.Sprintf("check existence of %s", resourceId.Describe())

    err := retryShort(action, nil, func() error {
        _, err := collection.GetYaml(resourceId.Name())
        return err
    })

    if err == ErrNotFound {
        return false, nil
    }

    return err == nil, nil
}

func (client *KubeClient) Delete(resourceId *KubeResourceId) error {
    if resourceId.CannotBeDeleted() {
        return nil
    }

    collection := resourceId.GetCollection(client.restClient)
    action := fmt.Sprintf("delete %s", resourceId.Describe())

    err := retryShort(action, nil, func() error {
        return collection.Delete(resourceId.Name())
    })

    if err == ErrNotFound {
        return nil
    }

    if err != nil {
        dumpErrorToFile(action, nil, err)
    }

    return err
}

func retryLong(action string, content []byte, do func() error) error {
    return retry(200, action, content, do) // 10 minutes
}

func retryShort(action string, content []byte, do func() error) error {
    return retry(3, action, content, do) // 3 times
}

func retry(n int, action string, content []byte, do func() error) (err error) {
    var done bool
    for i := 0; i < n; i++ {
        if i != 0 {
            time.Sleep(3 * time.Second)
        }

        if done, err = try(do); done {
            return
        }

        if err != nil {
            dumpError(os.Stderr, action, content, err)
        }
    }
    return
}

func try(do func() error) (bool, error) {
    err := do()
    if err == nil {
        return true, nil
    }

    restErr, ok := err.(*rest_error.Error)
    if !ok {
        return false, err
    }

    if restErr.IsClientError() {
        if restErr.Code == http.StatusNotFound {
            return true, ErrNotFound
        }
        if restErr.Code == http.StatusConflict {
            return true, ErrConflict
        }
        if restErr.Code == http.StatusForbidden { // Illegal Kubernetes state, need to retry
            return false, restErr
        }
        return true, restErr
    }

    return false, restErr
}

func dumpErrorToFile(action string, content []byte, err error) {
    file, err2 := os.Create("kubernetes-error.log")
    if err2 != nil {
        return
    }

    defer file.Close()
    defer file.Sync()

    fmt.Fprintf(file, "%s\n\n", time.Now().String())

    dumpError(file, action, content, err)
}

func dumpError(writer io.Writer, action string, content []byte, err error) {
    fmt.Fprintf(writer, "Failed to %s: %v\n\n", action, err)

    if content != nil {
        fmt.Fprintln(writer, "=============")
        fmt.Fprintln(writer, "== Content ==")
        fmt.Fprintln(writer, "=============")
        fmt.Fprintf(writer, "\n%s\n\n", content)
    }
}