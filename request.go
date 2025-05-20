package httpregistry

import (
	"encoding/json"
	"reflect"
	"regexp"
)

// Request represents a request that will be registered to a Registry to get matched against an incoming HTTP request.
// The match happens against the method, the headers and the URL interpreted as a regex
type Request struct {
	URL        string            `json:"url,omitempty"`
	Method     string            `json:"method,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       []byte            `json:"body,omitempty"`
	urlAsRegex regexp.Regexp
}

// WithURL returns a new request with the URL attribute set to URL
func (r Request) WithURL(URL string) Request {
	r.URL = URL
	r.urlAsRegex = *regexp.MustCompile(URL)
	return r
}

// WithMethod returns a new request with the method attribute set to method
func (r Request) WithMethod(method string) Request {
	r.Method = method
	return r
}

// WithHeader returns a new request with the header header set to value
func (r Request) WithHeader(header string, value string) Request {
	r.Headers[header] = value
	return r
}

// WithJSONHeader returns a new request with the header `Content-Type` set to `application/json`
func (r Request) WithJSONHeader() Request {
	r.Headers["Content-Type"] = "application/json"
	return r
}

// WithHeaders returns a new request with all the headers in headers applied.
// If multiple headers with the same name are defined only the last one is applied.
func (r Request) WithHeaders(headers map[string]string) Request {
	for k, v := range headers {
		r.Headers[k] = v
	}
	return r
}

// WithBody returns a new request with the method body set to body
func (r Request) WithBody(body []byte) Request {
	r.Body = body
	return r
}

// WithStringBody returns a new request with the method body set to body
func (r Request) WithStringBody(body string) Request {
	r.Body = []byte(body)
	return r
}

// WithJSONBody returns a new request with the method body set to the JSON encoded version of body and
// the Content-Type header set to "application/json".
// This method panics if body cannot be converted to JSON
func (r Request) WithJSONBody(body any) Request {
	r = r.WithJSONHeader()
	r.Body = mustMarshalJSON(body)
	return r
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

// NewRequest creates a new request designed to be registered to a Registry to get matched against an incoming HTTP request.
// This function is designed to be used in conjunction with other other receivers.
// For example
//
//	NewRequest().
//		WithURL("/users/1").
//		WithMethod(http.MethodPatch).
//		WithJSONHeader().
//		WithBody([]byte("{\"user\": \"John Schmidt\"}"))
func NewRequest() Request {
	r := Request{
		URL:        "",
		urlAsRegex: *regexp.MustCompile(".+"),
		Method:     "",
		Headers:    make(map[string]string),
		Body:       make([]byte, 0),
	}
	return r
}
