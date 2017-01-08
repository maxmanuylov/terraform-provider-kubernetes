package kubernetes_client

import (
    "fmt"
    "github.com/maxmanuylov/go-rest/client"
    "github.com/maxmanuylov/go-rest/error"
    "io"
    "net/http"
    "os"
    "time"
)

var (
    ErrNotFound = rest_error.NewByCode(http.StatusNotFound)
    ErrConflict = rest_error.NewByCode(http.StatusConflict)
)

type KubeClient struct {
    restClient *rest_client.Client
}

type errorHistory struct {
    error   error
    history []error
}

func (client *KubeClient) WaitForAPIServer() error {
    action := "connect to Kubernetes API server"

    eh := retryLong(action, nil, func() error {
        _, err := client.restClient.Do("GET", DefaultApiPath, rest_client.Json, nil)
        return err
    })

    if eh.error != nil {
        dumpErrorsToFile(action, nil, eh)
    }

    return eh.error
}

func (client *KubeClient) Create(resource *KubeResource) error {
    collection := resource.GetCollection(client.restClient)

    action := fmt.Sprintf("create %s", resource.Describe())
    content := resource.PrepareContent()

    eh := retryLong(action, content, func() error {
        _, err := collection.CreateYaml(content)
        if err == http.ErrNoLocation {
            return nil
        }
        return err
    })

    if eh.error == ErrConflict { // resource already exists
        return client.doUpdate(resource, content)
    }

    if eh.error != nil {
        dumpErrorsToFile(action, content, eh)
        return eh.error
    }

    return nil
}

func (client *KubeClient) Update(resource *KubeResource) error {
    return client.doUpdate(resource, resource.PrepareContent())
}

func (client *KubeClient) doUpdate(resource *KubeResource, content []byte) error {
    collection := resource.GetCollection(client.restClient)
    action := fmt.Sprintf("update %s", resource.Describe())

    eh := retryLong(action, content, func() error {
        return collection.ReplaceYaml(resource.Name(), content)
    })

    if eh.error != nil {
        dumpErrorsToFile(action, content, eh)
    }

    return eh.error
}

func (client *KubeClient) Exists(resourceId *KubeResourceId) (bool, error) {
    collection := resourceId.GetCollection(client.restClient)
    action := fmt.Sprintf("check existence of %s", resourceId.Describe())

    eh := retryShort(action, nil, func() error {
        _, err := collection.GetYaml(resourceId.Name())
        return err
    })

    if eh.error == ErrNotFound {
        return false, nil
    }

    return eh.error == nil, nil
}

func (client *KubeClient) Delete(resourceId *KubeResourceId) error {
    if resourceId.CannotBeDeleted() {
        return nil
    }

    collection := resourceId.GetCollection(client.restClient)
    action := fmt.Sprintf("delete %s", resourceId.Describe())

    eh := retryShort(action, nil, func() error {
        return collection.Delete(resourceId.Name())
    })

    if eh.error == ErrNotFound {
        return nil
    }

    if eh.error != nil {
        dumpErrorsToFile(action, nil, eh)
    }

    return eh.error
}

func retryLong(action string, content []byte, do func() error) *errorHistory {
    return retry(200, action, content, do) // 10 minutes
}

func retryShort(action string, content []byte, do func() error) *errorHistory {
    return retry(3, action, content, do) // 3 times
}

func retry(n int, action string, content []byte, do func() error) *errorHistory {
    eh := &errorHistory{
        history: make([]error, 0),
    }

    var done bool

    for i := 0; i < n; i++ {
        if i != 0 {
            time.Sleep(3 * time.Second)
        }

        done, eh.error = try(do)
        if eh.error != nil {
            eh.history = append(eh.history, eh.error)
        }

        if done {
            return eh
        }

        if eh.error != nil {
            dumpErrors(os.Stderr, action, content, eh.error)
        }
    }

    return eh
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

func dumpErrorsToFile(action string, content []byte, eh *errorHistory) {
    file, err := os.Create("kubernetes-error.log")
    if err != nil {
        return
    }

    defer file.Close()
    defer file.Sync()

    fmt.Fprintf(file, "%s\n\n", time.Now().String())

    dumpErrors(file, action, content, eh.history...)
}

func dumpErrors(writer io.Writer, action string, content []byte, errors... error) {
    fmt.Fprintf(writer, "Failed to %s\n\n", action)

    for _, err := range errors {
        fmt.Fprintln(writer, "=================================================================")
        fmt.Fprintf(writer, "\n%v\n\n", err)
    }

    if content != nil {
        fmt.Fprintln(writer, "=============")
        fmt.Fprintln(writer, "== Content ==")
        fmt.Fprintln(writer, "=============")
        fmt.Fprintf(writer, "\n%s\n\n", content)
    }
}
