package main

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"os"
)

func main() {
	err := execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func execute() error {
	ctx := context.Background()
	conf, err := config.NewFromCFHome()
	if err != nil {
		return err
	}
	conf.WithSkipTLSValidation(true)
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	// grab the first org we find
	org, err := cf.Organizations.First(ctx, nil)
	if err != nil {
		return err
	}

	// add a label and annotation
	fmt.Printf("adding metadata label and annotation to org %s\n", org.Name)
	m := &resource.Metadata{}
	m.SetLabel("", "example-label1", "short-label")
	m.SetLabel("cf.example.org", "example-label2", "prefixed-label")
	m.SetAnnotation("", "example-annotation1", "short-annotation")
	m.SetAnnotation("cf.example.org", "example-annotation2", "prefixed-annotation")
	orgUpdate := &resource.OrganizationUpdate{
		Metadata: m,
	}
	org, err = cf.Organizations.Update(ctx, org.GUID, orgUpdate)
	if err != nil {
		return err
	}

	fmt.Printf("org %s metadata:\n", org.Name)
	printMap(org.Metadata.Labels)
	printMap(org.Metadata.Annotations)

	// now clear out the metadata we added
	fmt.Printf("clearing metadata label and annotation from org %s\n", org.Name)
	orgUpdate.Metadata.Clear()
	org, err = cf.Organizations.Update(ctx, org.GUID, orgUpdate)
	if err != nil {
		return err
	}

	fmt.Printf("org %s metadata:\n", org.Name)
	printMap(org.Metadata.Labels)
	printMap(org.Metadata.Annotations)

	return nil
}

func printMap(m map[string]*string) {
	for k, v := range m {
		fmt.Printf("%s=%s\n", k, *v)
	}
}
