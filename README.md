# go-cfclient

### Overview

[![GoDoc](https://godoc.org/github.com/cloudfoundry-community/go-cfclient?status.png)](https://godoc.org/github.com/cloudfoundry-community/go-cfclient)

`cfclient` is a package to assist you in writing apps that need information out of [Cloud Foundry](http://cloudfoundry.org). It provides functions and structures to retrieve


### Usage

`go get github.com/cloudfoundry-community/go-cfclient`

```go
package main

import (
	"github.com/cloudfoundry-community/go-cfclient"
)

func main() {
  c := &Config{
    ApiAddress:   "https://api.10.244.0.34.xip.io",
    LoginAddress: "https://login.10.244.0.34.xip.io",
    Username:     "admin",
    Password:     "admin",
  }
  client := NewClient(c)
  apps := client.ListApps()
  fmt.Println(apps)
}
```

### Contributing

Pull requests welcomed. Please ensure you make your changes in a branch off of the `develop` branch, not the `master` branch.
