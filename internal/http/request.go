package http

import "io"

// Request is used to help build up an HTTP request
type Request struct {
	method        string
	pathAndQuery  string
	contentType   string
	contentLength *int64

	// can set one or the other but not both
	body   io.Reader
	object any

	// arbitrary headers
	headers map[string]string
}

func NewRequest(method, pathAndQuery string) *Request {
	return &Request{
		method:       method,
		pathAndQuery: pathAndQuery,
		headers:      make(map[string]string),
	}
}

func (r *Request) WithObject(obj any) *Request {
	r.object = obj

	// default content type to json if provided an object
	if r.contentType == "" {
		r.contentType = "application/json"
	}
	return r
}

func (r *Request) WithBody(body io.Reader) *Request {
	r.body = body

	// default content type to json if provided a body
	if r.contentType == "" {
		r.contentType = "application/json"
	}
	return r
}

func (r *Request) WithContentType(contentType string) *Request {
	r.contentType = contentType
	return r
}

func (r *Request) WithContentLength(len int64) *Request {
	r.contentLength = &len
	return r
}

func (r *Request) WithHeader(name, value string) *Request {
	r.headers[name] = value
	return r
}
