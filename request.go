package httpregistry

import (
	"encoding/json"
	"reflect"
	"regexp"
)

// Request represents a request that will be registered to a Registry to get matched against an incoming HTTP request.
// The match happens against the method, the headers and the URL interpreted as a regex
type Request struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers,omitempty"`
	urlAsRegex regexp.Regexp
}

// Equal checks if a request is identical to another
func (r Request) Equal(r2 Request) bool {
	return reflect.DeepEqual(r, r2)
}

// String returns a human readable representation
func (r Request) String() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic("cannot marshal request")
	}
	return string(bytes)
}

// RequestOption represents a option that can be passed to NewRequest when creating a new request.
// NewRequest uses the Option patters to make it easy to configure the request behavior.
// For example
//
//	NewRequest(GET, "/foo", WithResponseHeader("Content-Type", "application/json"))
//
// will return a request with the desired content type
type RequestOption func(*Request)

// WithRequestHeader allows to add any header to a request.
// If multiple headers with the same name are defined only the last one is applied.
// To define multiple headers it is recommended to use WithRequestHeaders but it is possible to chain multiple calls of WithResponseHeader.
func WithRequestHeader(header string, value string) RequestOption {
	return func(r *Request) {
		r.Headers[header] = value
	}
}

// WithRequestHeaders allows to add any number of headers to a request.
// If multiple headers with the same name are defined only the last one is applied.
func WithRequestHeaders(headers map[string]string) RequestOption {
	return func(r *Request) {
		for k, v := range headers {
			r.Headers[k] = v
		}
	}
}

// NewRequest creates a new request designed to be registered to a Registry to get matched against an incoming HTTP request.
// The match happens against the method and the URL interpreted as a regex.
// If further options are passed they will be applied one after the other.
func NewRequest(method string, url string, options ...RequestOption) Request {
	r := Request{
		URL:        url,
		urlAsRegex: *regexp.MustCompile(url),
		Method:     method,
		Headers:    make(map[string]string),
	}
	for _, o := range options {
		o(&r)
	}

	return r
}

// NewJSONRequest creates a new request designed to be registered to a Registry to get matched against an incoming HTTP request.
// The match happens against the method and the URL interpreted as a regex.
// If further options are passed they will be applied one after the other.
// The "Content-Type" header is set to "application/json"
func NewJSONRequest(method string, url string, options ...RequestOption) Request {
	r := Request{
		URL:        url,
		urlAsRegex: *regexp.MustCompile(url),
		Method:     method,
		Headers:    make(map[string]string),
	}
	for _, o := range options {
		o(&r)
	}
	r.Headers["Content-Type"] = "application/json"
	return r
}
