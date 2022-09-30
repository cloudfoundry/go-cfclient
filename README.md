# go-cfclient

[![build workflow](https://github.com/cloudfoundry-community/go-cfclient/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/cloudfoundry-community/go-cfclient/actions/workflows/build.yml)
[![GoDoc](https://godoc.org/github.com/cloudfoundry-community/go-cfclient?status.svg)](http://godoc.org/github.com/cloudfoundry-community/go-cfclient)
[![Report card](https://goreportcard.com/badge/github.com/cloudfoundry-community/go-cfclient)](https://goreportcard.com/report/github.com/cloudfoundry-community/go-cfclient)

## Overview

`go-cfclient` is a package to assist you in writing apps that need to interact the [Cloud Foundry](http://cloudfoundry.org)
Cloud Controller [v3 API](https://v3-apidocs.cloudfoundry.org). The v2 API is no longer supported, however if you still
need to use the older API you may use the v2 branch.

## Usage
It's recommended you use the latest tagged version of the library and upgrade to newer version at your convenience.
This library now follows semantic versioning of releases.
```
go get github.com/cloudfoundry-community/go-cfclient@v1.0.0
```

Some example code:

```go
package main

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient"
)

func main() {
	c := &cfclient.Config{
		ApiAddress: "https://api.10.244.0.34.xip.io",
		Username:   "admin",
		Password:   "secret",
	}
	client, _ := cfclient.NewClient(c)
	apps, _ := client.ListApps()
	fmt.Println(apps)
}
```

## Development

All development takes place on feature branches and is merged to the `main` branch. Therefore the main
branch is considered a potentially unstable branch until a new release (see below) is cut.

```shell
make all
```

Please attempt to use standard go naming conventions for all structs, for example prefer GUID over Guid. Packages
should roughly follow the v3 API resources although short names are preferred, for example:

|- app
|- space
|- org
|- route

### Releases

This library uses [semantic versioning](https://go.dev/doc/modules/version-numbers) to release new features,
bug fixes or other breaking changes [via git tags](https://go.dev/doc/modules/publishing).

### Errors

If the Cloud Foundry error definitions change at <https://github.com/cloudfoundry/cloud_controller_ng/blob/master/vendor/errors/v2.yml>
then the error predicate functions in this package need to be regenerated.

To do this, simply use Go to regenerate the code:

```shell
make generate
```

## Contributing

Pull requests welcome. Please ensure you run all the unit tests, go fmt the code, and golangci-lint via `make all`
