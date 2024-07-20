# Netbox BGP Plugin Terraform Provider

This provider allows managing BGP resources in Netbox, when the [Netbox BGP plugin](https://github.com/netbox-community/netbox-bgp) is installed.

The provider is intentionally using a similar structure to the [`e-breuninger/netbox`](https://registry.terraform.io/providers/e-breuninger/netbox/latest) provider. If you're already using that provider, this provider should work smoothly alongside it.

## Configure

Example configuration:

```tf
provider "netboxbgp" {
  server_url = "https://netbox.my-company.net"
  api_token  = var.netbox_api_token
}
```

You can also set the provider config from environment variables:

- `NETBOX_SERVER_URL` in place of `server_url`
- `NETBOX_API_TOKEN` in place of `api_token`

For more details and additional properties, see [the docs](./docs/index.md).

## Development

### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22
- [Docker](https://docs.docker.com/desktop/)

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

```shell
make testacc
```
