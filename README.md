# go-cfclient

## Overview
`cfclient` is a package to assist you in writing apps that need to interact with the [Cloud Foundry](http://cloudfoundry.org)
v2 cloud controller API.

## v2 go-cfclient deprecated
The v2 version of the client and corresponding v2 cloud controller (CC) API is deprecated. Please start using the v3 version of this
client and CC API. This v2 branch is only kept around to support critical bug fixes.

## Upgrading the v2 go-cfclient
If you're currently using an old version of the go-cfclient and need to upgrade to the latest version that still
supports the v2 CC API, then you'll need to go get the "new" v2 module.

```shell
$ go get -u github.com/cloudfoundry-community/go-cfclient/v2
```

Update your go import statements as necessary in your go source files, then finally:
```shell
$ go mod tidy
```

## Usage
```
go get github.com/cloudfoundry-community/go-cfclient/v2
```
Some example code:

```go
package main

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient/v2"
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

### Paging Results

The API supports paging results via query string parameters. All of the v3 ListV3*ByQuery functions support paging. Only a subset of v2 function calls support paging the results:

- ListSpacesByQuery
- ListOrgsByQuery
- ListAppsByQuery
- ListServiceInstancesByQuery
- ListUsersByQuery

You can iterate over the results page-by-page using a function similar to this one:

```go
func processSpacesOnePageAtATime(client *cfclient.Client) error {
	page := 1
	pageSize := 50

	q := url.Values{}
	q.Add("results-per-page", strconv.Itoa(pageSize))

	for {
		// get the current page of spaces
		q.Set("page", strconv.Itoa(page))
		spaces, err := client.ListSpacesByQuery(q)
		if err != nil {
			fmt.Printf("Error getting spaces by query: %s", err)
			return err
		}

		// do something with each space
		fmt.Printf("Page %d:\n", page)
		for _, s := range spaces {
			fmt.Println("  " + s.Name)
		}

		// if we hit an empty page or partial page, that means we're done
		if len(spaces) < pageSize {
			break
		}

		// next page
		page++
	}
	return nil
}
```

## Development

```shell
make all
```

### Errors

If the Cloud Foundry error definitions change at <https://github.com/cloudfoundry/cloud_controller_ng/blob/master/vendor/errors/v2.yml>
then the error predicate functions in this package need to be regenerated.

To do this, simply use Go to regenerate the code:

```shell
make generate
```

