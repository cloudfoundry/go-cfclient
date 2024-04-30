# go-cfclient

[![build workflow](https://github.com/cloudfoundry/go-cfclient/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/cloudfoundry/go-cfclient/actions/workflows/build.yml)
[![GoDoc](https://godoc.org/github.com/cloudfoundry/go-cfclient/v3?status.svg)](http://godoc.org/github.com/cloudfoundry/go-cfclient/v3)
[![Report card](https://goreportcard.com/badge/github.com/cloudfoundry/go-cfclient/v3)](https://goreportcard.com/report/github.com/cloudfoundry/go-cfclient/v3)

## Overview
`go-cfclient` is a go module library to assist you in writing apps that need to interact the [Cloud Foundry](http://cloudfoundry.org)
Cloud Controller [v3 API](https://v3-apidocs.cloudfoundry.org). The v2 API is no longer supported, however if you _really_ 
need to use the older API you may use the go-cfclient v2 branch and releases.

__NOTE__ - The v3 version in the main branch is currently under development and may have **breaking changes** until a v3.0.0 release is
 cut. Until then, you may want to pin to a specific v3.0.0-alpha.x release.

## Installation
go-cfclient is compatible with modern Go releases in module mode, with Go installed:
```
go get github.com/cloudfoundry/go-cfclient/v3
```
Will resolve and add the package to the current development module, along with its dependencies. Eventually this
library will cut releases that will be tagged with v3.0.0, v3.0.1 etc, see the Versioning section below.

## Usage
- [Authentication](./README.md#authentication)
- [Resources](./README.md#resources)
- [Pagination](./README.md#pagination)
- [Asynchronous Jobs](./README.md#asynchronous-jobs)
- [Error Handling](./README.md#error-handling)
- [Migrating v2 to v3](./README.md#migrating-v2-to-v3)

Using go modules import the client, config and resource packages:
```go
import (
    "github.com/cloudfoundry/go-cfclient/v3/client"
    "github.com/cloudfoundry/go-cfclient/v3/config"
    "github.com/cloudfoundry/go-cfclient/v3/resource"
)
```

### Authentication
Construct a new CF client configuration object. The configuration object configures how the client will authenticate to the 
CF API. There are various supported auth mechanisms.

The simplest being - use the existing CF CLI configuration and auth token:
```go
cfg, _ := config.NewFromCFHome()
cf, _ := client.New(cfg)
```
Username and password:
```go
cfg, _ := config.New("https://api.example.org", config.UserPassword("user", "pass"))
cf, _ := client.New(cfg)
```
Client and client secret:
```go
cfg, _ := config.New("https://api.example.org", config.ClientCredentials("cf", "secret"))
cf, _ := client.New(cfg)
```
Static OAuth token, which requires both an access and refresh token:
```go
cfg, _ := config.New("https://api.example.org", config.Token(accessToken, refreshToken))
cf, _ := client.New(cfg)
```
For more detailed examples of using the various authentication and configuration options, see the
[auth example](./examples/auth/main.go).

### Resources
The services of a client divide the API into logical chunks and correspond to the structure of the CF API documentation
at https://v3-apidocs.cloudfoundry.org. In other words each major resource type has its own service client that
is accessible via the main client instance.
```go
apps, _ := cf.Applications.ListAll(context.Background(), nil)
for _, app := range apps {
    fmt.Printf("Application %s is %s\n", app.Name, app.State)
}
```
All clients and their functions that interact with the CF API live in the `client` package. The client package
is responsible for making HTTP requests using the resources defined in the `resource` package. All generic serializable
resource definitions live in the `resource` package and could be reused with other client's outside this library.

__NOTE__ - Using the context package you can easily pass cancellation signals and deadlines to various client calls
for handling a request. In case there is no context available, then `context.Background()` can be used as a starting
point.

### Pagination
All requests for resource collections (apps, orgs, spaces etc) support pagination. Pagination options are described
in the client.ListOptions struct and passed to the list methods directly or as an embedded type of a more specific
list options struct (for example client.AppListOptions).

Example iterating through all apps one page at a time:
```go
opts := client.NewAppListOptions()
for {
    apps, pager, _ := cf.Applications.List(context.Background(), opts)
    for _, app := range apps {
        fmt.Printf("Application %s is %s\n", app.Name, app.State)
    }  
    if !pager.HasNextPage() {
        break
    }
    pager.NextPage(opts)
}
```
If you'd rather get _all_ of the resources in one go and not worry about paging, every collection
has a corresponding `All` method that gathers all the resources from every page before returning. While this may be
convenient, this could have negative performance consequences on larger foundations/collections.
```go
opts := client.NewAppListOptions()
apps, _ := cf.Applications.ListAll(context.Background(), opts)
for _, app := range apps {
    fmt.Printf("Application %s is %s\n", app.Name, app.State)
}
```

### Asynchronous Jobs
Some API calls are long-running so immediately return a JobID (GUID) instead of waiting and returning a resource. In
those cases you only know if the job was accepted. You will need to poll the Job API to find out when the job
finishes. There's a `PollComplete` utility function that you can use to block until the job finishes:
```go
jobGUID, err := cf.Manifests.ApplyManifest(context.Background(), spaceGUID, manifest))
if err != nil {
    return err
}
opts := client.NewPollingOptions()
err = cf.Jobs.PollComplete(context.Background(), jobGUID, opts)
if err != nil {
    return err
}
```
The timeout and polling interval can be configured using the PollingOptions struct.

The PollComplete function will return a nil error if the job completes successfully. If PollComplete
times out waiting for the job to complete a `client.AsyncProcessTimeoutError` is returned. If the job itself
failed then the job API is queried for the job error which is then returned as a `resource.CloudFoundryError`
which can be inspected to find the failure cause.

### Error Handling
All client methods will return a `resource.CloudFoundryError` or sub-type for any response that isn't a 200 level
status code. All CF errors have a corresponding error code and the client uses those codes to construct a specific
client side error type. This allows you to easily branch your logic based off specific API error codes using one of
the many `resource.IsSomeTypeOfError(err error)` functions, for example:
```go
params, err := cf.ServiceCredentialBindings.GetParameters(guid)
if resource.IsServiceFetchBindingParametersNotSupportedError(err) {
    fmt.Println(err.(resource.CloudFoundryError).Detail)
} else if err != nil {
    return err // all other errors
} else {
    fmt.Printf("Parameters: %v\n", params)
}
```

### Migrating v2 to v3
A very basic example using the v2 client:
```go
c := &cfclient.Config{
    ApiAddress: "https://api.sys.example.com",
    Username:   "user",
    Password:   "password",
}
client, _ := cfclient.NewClient(c)
apps, _ := client.ListApps()
for _, a := range apps {
	fmt.Println(a.Name)
}
```
Converted to do the same in the v3 client:
```go
cfg, _ := config.New("https://api.sys.example.org", config.UserPassword("user", "pass"))
client, _ := client.New(cfg)
apps, _ := client.Applications.ListAll(context.Background(), nil)
for _, app := range apps {
    fmt.Println(app.Name)
}
```
If you need to migrate over to the new client iteratively you can do that by referencing both the old and new modules
simultaneously and creating two separate client instances - one for each module version.

Some of the main differences between the old client and the new v3 client include:
- The old v2 client supported most v2 resources and few a v3 resources. The new v3 client supports all v3 resources and
no v2 resources. While most v2 resources are similar to their v3 counterparts, some code changes will need to be made.
- The v2 client had a single type that contained resource functions disambiguated by function name. The v3
client has a separate client nested under the main client for each resource. For example:
`client.ListApps()` vs `client.Applications.ListAll()`
- All v3 client functions take a cancellable context object as their first parameter.
- All v3 client list functions take a type specific options struct that support type safe filtering.
- All v3 client list functions support paging natively in the client for ease of use.
- The v3 client supports shared access across goroutines.

## Versioning
In general, go-cfclient follows [semver](https://go.dev/doc/modules/version-number) as closely as we can for [tagging
releases](https://go.dev/doc/modules/publishing) of the package. We've adopted the following versioning policy:

- We increment the major version with any incompatible change to non-preview functionality, including changes to the exported Go API surface or behavior of the API.
- We increment the minor version with any backwards-compatible changes to functionality
- We increment the patch version with any backwards-compatible bug fixes.

## Development

All development takes place on feature branches and is merged to the `main` branch. Therefore the main
branch is considered a potentially unstable branch until a new release (see below) is cut.

```shell
make all
```

Please attempt to use standard go naming conventions for all structs, for example use GUID over Guid. All client
functions should have at least once basic unit test.

### Errors

If the Cloud Foundry error definitions change at <https://github.com/cloudfoundry/cloud_controller_ng/blob/master/vendor/errors/v2.yml>
then the error predicate functions in this package need to be regenerated.

To do this, simply use Go to regenerate the code:

```shell
make generate
```

## Contributing

Pull requests welcome. Please ensure you run all the unit tests, go fmt the code, and golangci-lint via `make all`
