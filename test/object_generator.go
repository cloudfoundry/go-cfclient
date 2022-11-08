package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"text/template"
)

const defaultAPIResourcePath = "https://api.example.org/v3/somepagedresource"

type ResourceResult struct {
	Resource string

	// extra included resources
	// https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#resources-with-includes
	Apps             []string
	Spaces           []string
	Organizations    []string
	Domains          []string
	Users            []string
	ServiceOfferings []string
	ServiceInstances []string
	Routes           []string
}

type PagedResult struct {
	Resources []string

	// extra included resources
	// https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#resources-with-includes
	Apps             []string
	Spaces           []string
	Organizations    []string
	Domains          []string
	Users            []string
	ServiceOfferings []string
	ServiceInstances []string
	Routes           []string
}

type resourceTemplate struct {
	GUID string
	Name string
}

type resultTemplate struct {
	TotalResults int
	TotalPages   int
	FirstPage    string
	LastPage     string
	NextPage     string
	PreviousPage string

	Resources        string
	Apps             string
	Spaces           string
	Organizations    string
	Domains          string
	Users            string
	ServiceOfferings string
	ServiceInstances string
	Routes           string
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

func (o *ObjectJSONGenerator) AppFeature() string {
	r := resourceTemplate{}
	return o.template(r, "app_feature.json")
}

func (o *ObjectJSONGenerator) AppUsage() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "app_usage.json")
}

func (o *ObjectJSONGenerator) AuditEvent() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "audit_event.json")
}

func (o *ObjectJSONGenerator) AppUpdateEnvVars() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "app_update_envvar.json")
}

func (o *ObjectJSONGenerator) AppEnvironment() string {
	r := resourceTemplate{
		Name: RandomName(),
	}
	return o.template(r, "app_environment.json")
}

func (o *ObjectJSONGenerator) AppEnvVar() string {
	r := resourceTemplate{}
	return o.template(r, "app_envvar.json")
}

func (o *ObjectJSONGenerator) AppSSH() string {
	return o.template(resourceTemplate{}, "app_ssh.json")
}

func (o *ObjectJSONGenerator) AppPermission() string {
	return o.template(resourceTemplate{}, "app_permissions.json")
}

func (o *ObjectJSONGenerator) Build() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "build.json")
}

func (o *ObjectJSONGenerator) Buildpack() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "buildpack.json")
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

func (o *ObjectJSONGenerator) EnvVarGroup() string {
	return o.template(resourceTemplate{}, "environment_variable_group.json")
}

func (o *ObjectJSONGenerator) FeatureFlag() string {
	r := resourceTemplate{
		Name: RandomName(),
	}
	return o.template(r, "feature_flag.json")
}

func (o *ObjectJSONGenerator) IsolationSegment() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "isolation_segment.json")
}

func (o *ObjectJSONGenerator) IsolationSegmentRelationships() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "isolation_segment_relationships.json")
}

func (o *ObjectJSONGenerator) Job() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "job.json")
}

func (o *ObjectJSONGenerator) Manifest() string {
	r := resourceTemplate{}
	return o.template(r, "manifest.yml")
}

func (o *ObjectJSONGenerator) ManifestDiff() string {
	r := resourceTemplate{}
	return o.template(r, "manifest_diff.yml")
}

func (o *ObjectJSONGenerator) Organization() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "org.json")
}

func (o *ObjectJSONGenerator) OrganizationUsageSummary() string {
	r := resourceTemplate{}
	return o.template(r, "org_usage_summary.json")
}

func (o *ObjectJSONGenerator) OrganizationQuota() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "org_quota.json")
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

func (o *ObjectJSONGenerator) Process() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "process.json")
}

func (o *ObjectJSONGenerator) ProcessStats() string {
	r := resourceTemplate{}
	return o.template(r, "process_stats.json")
}

func (o *ObjectJSONGenerator) ResourceMatch() string {
	r := resourceTemplate{}
	return o.template(r, "resource_match.json")
}

func (o *ObjectJSONGenerator) Revision() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
	}
	return o.template(r, "revision.json")
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

func (o *ObjectJSONGenerator) ServiceBroker() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "service_broker.json")
}

func (o *ObjectJSONGenerator) SecurityGroup() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "security_group.json")
}

func (o *ObjectJSONGenerator) ServiceCredentialBinding() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "service_credential_binding.json")
}

func (o *ObjectJSONGenerator) ServiceInstance() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "service_instance.json")
}

func (o *ObjectJSONGenerator) ServiceOffering() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "service_offering.json")
}

func (o *ObjectJSONGenerator) Space() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "space.json")
}

func (o *ObjectJSONGenerator) User() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "user.json")
}

func (o *ObjectJSONGenerator) Stack() string {
	r := resourceTemplate{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.template(r, "stack.json")
}

// ResourceWithInclude merges the included resources under the primary resource's included key
func (o ObjectJSONGenerator) ResourceWithInclude(rr ResourceResult) []string {
	j := map[string]any{}
	err := json.Unmarshal([]byte(rr.Resource), &j)
	if err != nil {
		panic(err)
	}

	t, err := template.New("res").Parse(singleTemplate)
	if err != nil {
		panic(err)
	}
	p := resultTemplate{
		Apps:             strings.Join(rr.Apps, ","),
		Spaces:           strings.Join(rr.Spaces, ","),
		Organizations:    strings.Join(rr.Organizations, ","),
		Domains:          strings.Join(rr.Domains, ","),
		Users:            strings.Join(rr.Users, ","),
		Routes:           strings.Join(rr.Routes, ","),
		ServiceOfferings: strings.Join(rr.ServiceOfferings, ","),
		ServiceInstances: strings.Join(rr.ServiceInstances, ","),
	}

	var h bytes.Buffer
	err = t.Execute(&h, p)
	if err != nil {
		panic(err)
	}
	s := h.String()
	j["included"] = json.RawMessage(s)

	b, err := json.Marshal(&j)
	if err != nil {
		panic(err)
	}
	s = string(b)
	return []string{s}
}

// PagedWithInclude takes the list of resources and inserts them into a paged API response
func (o ObjectJSONGenerator) PagedWithInclude(pagesOfResourcesJSON ...PagedResult) []string {
	totalPages := len(pagesOfResourcesJSON)
	totalResults := 0
	for _, pageOfResourcesJSON := range pagesOfResourcesJSON {
		totalResults += len(pageOfResourcesJSON.Resources)
	}

	// iterate through each page of resources and build a list of paginated responses
	var resultPages []string
	for i, pageOfResourcesJSON := range pagesOfResourcesJSON {
		pageIndex := i + 1
		resourcesPerPage := len(pageOfResourcesJSON.Resources)

		p := resultTemplate{
			TotalResults:     totalResults,
			TotalPages:       totalPages,
			FirstPage:        fmt.Sprintf("%s?page=1&per_page=%d", defaultAPIResourcePath, resourcesPerPage),
			LastPage:         fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, totalPages, resourcesPerPage),
			Resources:        strings.Join(pageOfResourcesJSON.Resources, ","),
			Apps:             strings.Join(pageOfResourcesJSON.Apps, ","),
			Spaces:           strings.Join(pageOfResourcesJSON.Spaces, ","),
			Organizations:    strings.Join(pageOfResourcesJSON.Organizations, ","),
			Domains:          strings.Join(pageOfResourcesJSON.Domains, ","),
			Users:            strings.Join(pageOfResourcesJSON.Users, ","),
			Routes:           strings.Join(pageOfResourcesJSON.Routes, ","),
			ServiceOfferings: strings.Join(pageOfResourcesJSON.ServiceOfferings, ","),
			ServiceInstances: strings.Join(pageOfResourcesJSON.ServiceInstances, ","),
		}
		if pageIndex < totalPages {
			p.NextPage = fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, pageIndex+1, resourcesPerPage)
		}
		if pageIndex > 1 {
			p.PreviousPage = fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, pageIndex-1, resourcesPerPage)
		}

		t, err := template.New("page").Parse(listTemplate)
		if err != nil {
			panic(err)
		}
		var h bytes.Buffer
		err = t.Execute(&h, p)
		if err != nil {
			panic(err)
		}
		s := h.String()
		resultPages = append(resultPages, s)

	}
	return resultPages
}

func (o ObjectJSONGenerator) Paged(pagesOfResourcesJSON ...[]string) []string {
	var pagedResults []PagedResult
	for _, pageOfResourcesJSON := range pagesOfResourcesJSON {
		p := PagedResult{
			Resources: pageOfResourcesJSON,
		}
		pagedResults = append(pagedResults, p)
	}
	return o.PagedWithInclude(pagedResults...)
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

const listTemplate = `
{
  "pagination": {
    "total_results": {{.TotalResults}},
    "total_pages": {{.TotalPages}},
    "first": { "href": "{{.FirstPage}}" },
    "last": { "href": "{{.LastPage}}" },
    {{if .NextPage}}"next": { "href": "{{.NextPage}}" },{{else}}"next": null,{{end}}
    {{if .PreviousPage}}"previous": { "href": "{{.PreviousPage}}" }{{else}}"previous": null{{end}}
  },
  "resources": [
    {{.Resources}}
  ],
  "included": {
    "apps": [
      {{.Apps}}
    ],
    "spaces": [
      {{.Spaces}}
    ],
    "domains": [
      {{.Domains}}
    ],
    "users": [
      {{.Users}}
    ],
    "routes": [
      {{.Routes}}
    ],
    "service_offerings": [
      {{.ServiceOfferings}}
    ],
    "service_instances": [
      {{.ServiceInstances}}
    ],
    "organizations": [
      {{.Organizations}}
    ]
  }
}
`

const singleTemplate = `
{
    "apps": [
      {{.Apps}}
    ],
    "spaces": [
      {{.Spaces}}
    ],
    "domains": [
      {{.Domains}}
    ],
    "users": [
      {{.Users}}
    ],
    "routes": [
      {{.Routes}}
    ],
    "service_offerings": [
      {{.ServiceOfferings}}
    ],
    "service_instances": [
      {{.ServiceInstances}}
    ],
    "organizations": [
      {{.Organizations}}
    ]
}
`
