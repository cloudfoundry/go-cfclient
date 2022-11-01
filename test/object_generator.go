package test

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"text/template"
)

const defaultAPIResourcePath = "https://api.example.org/v3/somepagedresource"

type resourceTemplate struct {
	GUID string
	Name string
}

type paginationTemplate struct {
	TotalResults int
	TotalPages   int
	FirstPage    string
	LastPage     string
	NextPage     string
	PreviousPage string
}

type ObjectJSONGenerator struct {
}

func NewObjectJSONGenerator(seed int) *ObjectJSONGenerator {
	rand.Seed(int64(seed)) // stable random
	return &ObjectJSONGenerator{}
}

func (o *ObjectJSONGenerator) Application() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "app.json")
}

func (o *ObjectJSONGenerator) AppUpdateEnvVars() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "app_update_envvar.json")
}

func (o *ObjectJSONGenerator) AppEnvVars() string {
	r := resourceTemplate{
		Name: RandomName(),
	}
	return o.template(r, "app_envvar.json")
}

func (o *ObjectJSONGenerator) AppSSH() string {
	return o.template(resourceTemplate{}, "app_ssh.json")
}

func (o *ObjectJSONGenerator) AppPermission() string {
	return o.template(resourceTemplate{}, "app_permissions.json")
}

func (o *ObjectJSONGenerator) Droplet() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "droplet.json")
}

func (o *ObjectJSONGenerator) DropletAssociation() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "droplet_association.json")
}

func (o *ObjectJSONGenerator) Build() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "build.json")
}

func (o *ObjectJSONGenerator) Deployment() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "deployment.json")
}

func (o *ObjectJSONGenerator) Domain() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "domain.json")
}

func (o *ObjectJSONGenerator) DomainShared() string {
	return o.template(resourceTemplate{}, "domain_shared.json")
}

func (o *ObjectJSONGenerator) Organization() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "org.json")
}

func (o *ObjectJSONGenerator) Package() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "package.json")
}

func (o *ObjectJSONGenerator) PackageDocker() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "package_docker.json")
}

func (o *ObjectJSONGenerator) Role() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "role.json")
}

func (o *ObjectJSONGenerator) Route() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "route.json")
}

func (o *ObjectJSONGenerator) SecurityGroup() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "security_group.json")
}

func (o ObjectJSONGenerator) Paged(pagesOfResourcesJSON ...[]string) []string {
	totalPages := len(pagesOfResourcesJSON)
	totalResults := 0
	for _, pageOfResourcesJSON := range pagesOfResourcesJSON {
		totalResults += len(pageOfResourcesJSON)
	}

	// iterate through each page of resources and build a list of paginated responses
	var resultPages []string
	for i, pageOfResourcesJSON := range pagesOfResourcesJSON {
		pageIndex := i + 1
		resourcesPerPage := len(pageOfResourcesJSON)

		p := paginationTemplate{
			TotalResults: totalResults,
			TotalPages:   totalPages,
			FirstPage:    fmt.Sprintf("%s?page=1&per_page=%d", defaultAPIResourcePath, resourcesPerPage),
			LastPage:     fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, totalPages, resourcesPerPage),
		}
		if pageIndex < totalPages {
			p.NextPage = fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, pageIndex+1, resourcesPerPage)
		}
		if pageIndex > 1 {
			p.PreviousPage = fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, pageIndex-1, resourcesPerPage)
		}

		t, err := template.New("page").Parse(responseListHeader)
		if err != nil {
			panic(err)
		}
		var h bytes.Buffer
		err = t.Execute(&h, p)
		if err != nil {
			panic(err)
		}

		s := h.String() + strings.Join(pageOfResourcesJSON, ",") + responseListFooter
		resultPages = append(resultPages, s)

	}
	return resultPages
}

func (o ObjectJSONGenerator) Array(resourcesJSON ...string) string {
	return "[" + strings.Join(resourcesJSON, ",") + "]"
}

func (o *ObjectJSONGenerator) template(rt resourceTemplate, fileName string) string {
	p := path.Join("../test/template", fileName)
	f, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}

	t, err := template.New("resource").Parse(string(f))
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	err = t.Execute(&b, rt)
	if err != nil {
		panic(err)
	}
	return b.String()
}

const responseListHeader = `{
  "pagination": {
    "total_results": {{.TotalResults}},
    "total_pages": {{.TotalPages}},
    "first": { "href": "{{.FirstPage}}" },
    "last": { "href": "{{.LastPage}}" },
    {{if .NextPage}}"next": { "href": "{{.NextPage}}" },{{else}}"next": null,{{end}}
    {{if .PreviousPage}}"previous": { "href": "{{.PreviousPage}}" }{{else}}"previous": null{{end}}
  },
  "resources": [`

const responseListFooter = `]}`
