# Terraform Kubernetes Provider

This is a plugin for HashiCorp [Terraform](https://terraform.io), which helps deploying Kubernetes resources like pods, services, replication controllers, etc.

## Usage

- Download the plugin from [Releases](https://github.com/maxmanuylov/terraform-provider-kubernetes/releases) page.
- [Install](https://terraform.io/docs/plugins/basics.html) it, or put into a directory with configuration files.
- Create a sample configuration file `example.tf`:
```hcl
provider "k8s" {
  # Either "k8s" provider or "k8s_cluster" resource should be configured

  # Kubernetes API server, both HTTP and HTTPS are supported
  api_server = "https://192.168.0.1:6443"

  # TLS options are optional, see "tls" Terraform provider (built-in) for certificates/keys generating
  ca_cert = "<CA certificate content (PEM)>"
  client_cert = "<client certificate content (PEM)>"
  client_key = "<client private key content (PEM)>"
}

resource "k8s_cluster" "main" {
  # Either "k8s" provider or "k8s_cluster" resource should be configured

  # Kubernetes API server, both HTTP and HTTPS are supported
  api_server = "https://192.168.0.1:6443"

  # TLS options are optional, see "tls" Terraform provider (built-in) for certificates/keys generating
  ca_cert = "<CA certificate content (PEM)>"
  client_cert = "<client certificate content (PEM)>"
  client_key = "<client private key content (PEM)>"
}

resource "k8s_resource" "mypod" {
  # Optional; if specified, must link on the corresponding "k8s_cluster" resource; otherwise provider configuration is used
  cluster = "${k8s_cluster.main.cluster}"

  # Required; resource contents must be in JSON or YAML format
  contents = "${file("mypod.yaml")}"

  # Optional; specifies "contents" format; possible values are "yaml" (default) and "json"
  encoding = "yaml"
}
```
- Run:
```
$ terraform apply
```
