package client

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type OrgClient commonClient

const OrgsPath = "/v3/organizations"

type OrgListOptions struct {
	*ListOptions

	GUIDs Filter
	Names Filter
}

func NewOrgListOptions() *OrgListOptions {
	return &OrgListOptions{
		ListOptions: NewListOptions(),
	}
}

func (a OrgListOptions) ToQuerystring() url.Values {
	v := a.ListOptions.ToQueryString()
	v = appendQueryStrings(v, a.GUIDs.ToQueryString(GUIDsField))
	v = appendQueryStrings(v, a.Names.ToQueryString(NamesField))
	return v
}

func (o *OrgClient) Create(r resource.CreateOrganizationRequest) (*resource.Organization, error) {
	params := map[string]interface{}{
		"name": r.Name,
	}
	if r.Suspended != nil {
		params["suspended"] = r.Suspended
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	var org resource.Organization
	err := o.client.post(r.Name, OrgsPath, params, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (o *OrgClient) Delete(guid string) error {
	return o.client.delete(joinPath(OrgsPath, guid))
}

func (o *OrgClient) Get(guid string) (*resource.Organization, error) {
	var org resource.Organization
	err := o.client.get(joinPath(OrgsPath, guid), &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (o *OrgClient) ListAll() ([]*resource.Organization, error) {
	opts := NewOrgListOptions()
	var allOrgs []*resource.Organization
	for {
		orgs, pager, err := o.List(opts)
		if err != nil {
			return nil, err
		}
		allOrgs = append(allOrgs, orgs...)
		if !pager.NextPage(opts.ListOptions) {
			break
		}
	}
	return allOrgs, nil
}

func (o *OrgClient) List(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
	var res resource.ListOrganizationsResponse
	err := o.client.get(joinPathAndQS(opts.ToQuerystring(), OrgsPath), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := &Pager{
		pagination: res.Pagination,
	}
	return res.Resources, pager, nil
}

func (o *OrgClient) Update(guid string, r resource.UpdateOrganizationRequest) (*resource.Organization, error) {
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Suspended != nil {
		params["suspended"] = r.Suspended
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	var org resource.Organization
	err := o.client.patch(joinPath(OrgsPath, guid), params, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
