package client

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

type Pager struct {
	NextPageURL     string
	PreviousPageURL string

	nextPageQSReader     *querystringReader
	previousPageQSReader *querystringReader
}

func NewPager(pagination resource.Pagination) *Pager {
	return &Pager{
		NextPageURL:     pagination.Next.Href,
		PreviousPageURL: pagination.Previous.Href,
	}
}

func (p *Pager) HasNextPage() bool {
	q, err := newQuerystringReader(p.NextPageURL)
	if err != nil {
		return false
	}
	p.nextPageQSReader = q
	return true
}

func (p Pager) NextPage(opts *ListOptions) *ListOptions {
	if !p.HasNextPage() {
		return opts
	}
	opts.Page = p.nextPageQSReader.Int(PageField)
	opts.PerPage = p.nextPageQSReader.Int(PerPageField)
	return opts
}

func (p *Pager) HasPreviousPage() bool {
	q, err := newQuerystringReader(p.PreviousPageURL)
	if err != nil {
		return false
	}
	p.previousPageQSReader = q
	return true
}

func (p Pager) PreviousPage(opts *ListOptions) *ListOptions {
	if !p.HasPreviousPage() {
		return opts
	}
	opts.Page = p.previousPageQSReader.Int(PageField)
	opts.PerPage = p.previousPageQSReader.Int(PerPageField)
	return opts
}

type querystringReader struct {
	qs url.Values
}

func newQuerystringReader(pageURL string) (*querystringReader, error) {
	if pageURL == "" {
		return nil, errors.New("cannot parse an empty pageURL")
	}
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	return &querystringReader{
		qs: u.Query(),
	}, nil
}

func (r querystringReader) String(key string) string {
	return r.qs.Get(key)
}

func (r querystringReader) Int(key string) int {
	i, _ := strconv.Atoi(r.qs.Get(key))
	return i
}
