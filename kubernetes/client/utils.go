package kubernetes_client

import (
    "fmt"
    "github.com/maxmanuylov/go-rest/client"
    "github.com/maxmanuylov/go-rest/error"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/cluster"
    "github.com/maxmanuylov/terraform-provider-kubernetes/kubernetes/model"
    "github.com/maxmanuylov/utils/http/transport/tls"
    "io"
    "net/http"
    "os"
    "time"
)

func newTransport(cluster *kubernetes_cluster.Cluster) (http.RoundTripper, error) {
    if cluster.CaCert != "" && cluster.ClientCert != "" && cluster.ClientKey != "" {
        return tls_transport.New([]byte(cluster.CaCert), []byte(cluster.ClientCert), []byte(cluster.ClientKey))
    }
    return http.DefaultTransport, nil
}

type errorHistory struct {
    error   error
    history []error
}

func errInvalidEncoding(encoding string) error {
    return fmt.Errorf("Invalid encoding: %s", encoding)
}

func createResource(collection rest_client.Collection, encoding string, contents []byte) error {
    if encoding == kubernetes_model.EncodingJson {
        _, err := collection.CreateJson(contents)
        return err
    } else if encoding == kubernetes_model.EncodingYaml {
        _, err := collection.CreateYaml(contents)
        return err
    }
    return errInvalidEncoding(encoding)
}

func updateResource(collection rest_client.Collection, name, encoding string, contents []byte) error {
    if encoding == kubernetes_model.EncodingJson {
        return collection.ReplaceJson(name, contents)
    } else if encoding == kubernetes_model.EncodingYaml {
        return collection.ReplaceYaml(name, contents)
    }
    return errInvalidEncoding(encoding)
}

func retryLong(action string, contents []byte, do func() error) *errorHistory {
    return retry(200, action, contents, do) // 10 minutes
}

func retryShort(action string, contents []byte, do func() error) *errorHistory {
    return retry(3, action, contents, do) // 3 times
}

func retry(n int, action string, contents []byte, do func() error) *errorHistory {
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
            dumpErrors(os.Stderr, action, contents, eh.error)
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

func dumpErrorsToFile(action string, contents []byte, eh *errorHistory) {
    file, err := os.Create("kubernetes-error.log")
    if err != nil {
        return
    }

    defer file.Close()
    defer file.Sync()

    fmt.Fprintf(file, "%s\n\n", time.Now().String())

    dumpErrors(file, action, contents, eh.history...)
}

func dumpErrors(writer io.Writer, action string, contents []byte, errors... error) {
    fmt.Fprintf(writer, "Failed to %s\n\n", action)

    for _, err := range errors {
        fmt.Fprintln(writer, "=================================================================")
        fmt.Fprintf(writer, "\n%v\n\n", err)
    }

    if contents != nil {
        fmt.Fprintln(writer, "==============")
        fmt.Fprintln(writer, "== Contents ==")
        fmt.Fprintln(writer, "==============")
        fmt.Fprintf(writer, "\n%s\n\n", contents)
    }
}
