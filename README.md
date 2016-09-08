# Terraform Kubernetes Provider

This is a plugin for HashiCorp [Terraform](https://terraform.io), which helps deploying Kubernetes resources like pods, services, replication controllers, etc.

## Usage

- Download the plugin from [Releases](https://github.com/maxmanuylov/terraform-provider-kubernetes/releases) page.
- [Install](https://terraform.io/docs/plugins/basics.html) it, or put into a directory with configuration files.
- Create a sample configuration file `terraform.tf`:
```
resource "kubernetes_cluster" "main" {
  # Required, both HTTP and HTTPS are supported
  api_server = "https://192.168.0.1:6443"

  # TLS options are optional, see "tls" Terraform provider (built-in) for certificates/keys generating
  ca_cert = "<CA certificate content (PEM)>"
  client_cert = "<client certificate content (PEM)>"
  client_key = "<client private key content (PEM)>"
}

resource "kubernetes_resource" "mypod" {
  # Required, must link on the corresponding "kubernetes_cluster" resource
  cluster = "${kubernetes_cluster.main.cluster}"

  # Optional, default is "api/v1"
  api_path = "api/v1"

  # Optional, default is "default"
  namespace = "default"

  # Required
  collection = "pods"
  name = "mypod"

  # Optional
  labels {
    a = "b"
  }

  # Optional
  annotations {
    a = "b"
  }

  # Required, resource content must be in .yaml format and must NOT contain "apiVersion", "kind" and "metadata" sections
  content = "${file("mypod.yaml")}"
}
```
- Run:
```
$ terraform apply
```
