# Terraform Provider RedisDB

The Terraform RedisDB Provider is a plugin for Terraform that allows for management of data in a redis instance.

This provider can't create redis instances, but is instead designed to connect to them and read/write data. This allows reading data
from redis that might be useful in other Terraform configuration, or writing data to redis that was generated by other Terrafom resources.


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

### Terraform 0.13 and above

You can use the provider via the [Terraform provider registry](https://registry.terraform.io/providers/cruglobal/redisdb).

### Terraform 0.12 or manual installation

You can download a pre-built binary from the [releases](https://github.com/CruGlobal/terraform-provider-redisdb/releases) page, these are built using [goreleaser](https://goreleaser.com/) (the [configuration](.goreleaser.yml) is in the repo).

If you want to build from source, you can simply use `go build` in the root of the repository.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make install`. This will build the provider and put the provider binary in the `~/.terraform.d/plugins` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`. Running the test suite requires Docker to run a temporary redis-server. If tests fail, you will need to manually stop the redis docker container. `./redis.sh stop`

```sh
$ make testacc
```
