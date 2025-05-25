package httpregistry

import (
	"reflect"
	"regexp"
)

// DefaultRequest represents the request that is used when no request is specified.
// It will match any request.
// This is useful in combination with registry.GetMatchesForRequest so that it is possible to retrieve all the matches associated with it
var DefaultRequest = newRequestWithName("httpregistry.DefaultRequest")

// Request represents a request that will be registered to a Registry to get matched against an incoming HTTP request.
// The match happens against the method, the headers and the URL interpreted as a regex
type Request struct {
	name       string
	url        string
	method     string
	headers    map[string]string
	body       []byte
	urlAsRegex regexp.Regexp
}

// Equal checks if a request is identical to another
func (r Request) Equal(r2 Request) bool {
	return reflect.DeepEqual(r.url, r2.url) &&
		reflect.DeepEqual(r.method, r2.method) &&
		reflect.DeepEqual(r.headers, r2.headers) &&
		reflect.DeepEqual(r.body, r2.body) &&
		reflect.DeepEqual(r.urlAsRegex, r2.urlAsRegex)
}

// String returns the name associated with the request
func (r Request) String() string {
	return r.name
}

// WithName allows to add a name to a Request so that it can be better identified when debugging.
// By the default Request gets a sequential name that can be hard to identify if there are many of them.
// So if clarity is needed we recommend to change the default name.
func (r Request) WithName(name string) Request {
	r.name = name
	return r
}

// WithURL returns a new request with the URL attribute set to URL
func (r Request) WithURL(URL string) Request {
	r.url = URL
	r.urlAsRegex = *regexp.MustCompile(URL)
	return r
}

// WithMethod returns a new request with the method attribute set to method
func (r Request) WithMethod(method string) Request {
	r.method = method
	return r
}

// WithHeader returns a new request with the header header set to value
func (r Request) WithHeader(header string, value string) Request {
	r.headers[header] = value
	return r
}

// WithJSONHeader returns a new request with the header `Content-Type` set to `application/json`
func (r Request) WithJSONHeader() Request {
	r.headers["Content-Type"] = "application/json"
	return r
}

// WithHeaders returns a new request with all the headers in headers applied.
// If multiple headers with the same name are defined only the last one is applied.
func (r Request) WithHeaders(headers map[string]string) Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

// WithBody returns a new request with the method body set to body
func (r Request) WithBody(body []byte) Request {
	r.body = body
	return r
}

// WithStringBody returns a new request with the method body set to body
func (r Request) WithStringBody(body string) Request {
	r.body = []byte(body)
	return r
}

// WithJSONBody returns a new request with the method body set to the JSON encoded version of body and
// the Content-Type header set to "application/json".
// This method panics if body cannot be converted to JSON
func (r Request) WithJSONBody(body any) Request {
	r = r.WithJSONHeader()
	r.body = mustMarshalJSON(body)
	return r
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
	return newRequestWithName("")
}

// newRequestWithName creates a new request designed to be registered to a Registry to get matched against an incoming HTTP request.
// This function allows to set the name while creating the request,
// this has the advantage of not increasing the counter in the default naming schema.
//
// This function is designed to be used in conjunction with other other receivers.
// For example
//
//	newRequestWithName().
//		WithURL("/users/1").
//		WithMethod(http.MethodPatch).
//		WithJSONHeader().
//		WithBody([]byte("{\"user\": \"John Schmidt\"}"))
func newRequestWithName(name string) Request {
	r := Request{
		name:       name,
		url:        "",
		urlAsRegex: *regexp.MustCompile(".+"),
		method:     "",
		headers:    make(map[string]string),
		body:       make([]byte, 0),
	}
	return r
}
