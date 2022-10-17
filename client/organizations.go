package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type OrgClient commonClient

const OrgsPath = "/v3/organizations"

type OrgListOptions struct {
	*ListOptions

	GUIDs Filter `filter:"guids,omitempty"`
	Names Filter `filter:"names,omitempty"`
}

func NewOrgListOptions() *OrgListOptions {
	return &OrgListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o *OrgClient) Create(r *resource.OrganizationCreate) (*resource.Organization, error) {
	var org resource.Organization
	err := o.client.post(r.Name, OrgsPath, r, &org)
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

func (o *OrgClient) List(opts *OrgListOptions) ([]*resource.Organization, *Pager, error) {
	var res resource.OrganizationList
	err := o.client.get(joinPathAndQS(opts.ToQueryString(opts), OrgsPath), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
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

func (o *OrgClient) Update(guid string, r *resource.OrganizationUpdate) (*resource.Organization, error) {
	var org resource.Organization
	err := o.client.patch(joinPath(OrgsPath, guid), r, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
