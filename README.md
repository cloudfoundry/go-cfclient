# go-cfclient

[![build workflow](https://github.com/cloudfoundry-community/go-cfclient/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/cloudfoundry-community/go-cfclient/actions/workflows/build.yml)
[![GoDoc](https://godoc.org/github.com/cloudfoundry-community/go-cfclient/v3?status.svg)](http://godoc.org/github.com/cloudfoundry-community/go-cfclient/v3)
[![Report card](https://goreportcard.com/badge/github.com/cloudfoundry-community/go-cfclient/v3)](https://goreportcard.com/report/github.com/cloudfoundry-community/go-cfclient/v3)

## Overview
`go-cfclient` is a go module library to assist you in writing apps that need to interact the [Cloud Foundry](http://cloudfoundry.org)
Cloud Controller [v3 API](https://v3-apidocs.cloudfoundry.org). The v2 API is no longer supported, however if you _really_ 
need to use the older API you may use the go-cfclient v2 branch and releases.

__NOTE__ - The v3 version in the main branch is currently under development and may have **breaking changes** until a v3.0.0 release is
 cut. Until then, you may want to pin to a specific v3.0.0-alpha.x release.

## Installation
go-cfclient is compatible with modern Go releases in module mode, with Go installed:
```
go get github.com/cloudfoundry-community/go-cfclient/v3
```
Will resolve and add the package to the current development module, along with its dependencies. Eventually this
library will cut releases that will be tagged with v3.0.0, v3.0.1 etc, see the Versioning section below.

## Usage
Using go modules import the client, config and resource packages:
```go
import (
    "github.com/cloudfoundry-community/go-cfclient/v3/client"
    "github.com/cloudfoundry-community/go-cfclient/v3/config"
    "github.com/cloudfoundry-community/go-cfclient/v3/resource"
)
```

### Authentication
Construct a new CF client configuration object. The configuration object configures how the client will authenticate to the 
CF API. There are various supported auth mechanisms, with the simplest being - use the existing CF CLI configuration and
auth token:
```go
cfg, _ := config.NewFromCFHome()
cf, _ := client.New(cfg)
```
You may also use username/password
```go
cfg, _ := config.NewUserPassword("https://api.example.org", "user", "pass")
cf, _ := client.New(cfg)
```
There is also client/secret and token config support.

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
If you'd rather have your code get _all_ of the resources in one go and not worry about paging, every collection
has a corresponding `All` method that gathers all the resources from every page before returning.
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
finishes. There's a `PollComplete` method to block until the job finishes:
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

The PollComplete method will return a nil err if the job completes successfully. If PollComplete
times out waiting for the job to complete a `client.AsyncProcessTimeoutError` is returned. If the job itself
failed then the job API is queried for the job error which is then returned as a `resource.CloudFoundryError`
which can be inspected to find the failure cause.

### Error handling
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
