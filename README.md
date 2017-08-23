# go-cfclient

### Overview

[![Build Status](https://img.shields.io/travis/cloudfoundry-community/go-cfclient.svg)](https://travis-ci.org/cloudfoundry-community/go-cfclient) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/cloudfoundry-community/go-cfclient)

`cfclient` is a package to assist you in writing apps that need information out of [Cloud Foundry](http://cloudfoundry.org). It provides functions and structures to retrieve


### Usage

```
go get github.com/cloudfoundry-community/go-cfclient
```

NOTE: Currently this project is not versioning its releases and so breaking changes might be introduced. Whilst hopefully notifications of breaking changes are made via commit messages, ideally your project will use a local vendoring system to lock in a version of `go-cfclient` that is known to work for you. This will allow you to control the timing and maintenance of upgrades to newer versions of this library.

Some example code:

```go
package main

import (
	"github.com/cloudfoundry-community/go-cfclient"
)

func main() {
  c := &cfclient.Config{
    ApiAddress:   "https://api.10.244.0.34.xip.io",
    Username:     "admin",
    Password:     "admin",
  }
  client, _ := cfclient.NewClient(c)
  apps, _ := client.ListApps()
  fmt.Println(apps)
}
```

### Developing & Contributing

You can use Godep to restore the dependency
Tested with go1.5.3
```bash
godep go build
```

Pull requests welcome.
