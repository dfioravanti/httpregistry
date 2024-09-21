package httpregistry

import (
	"encoding/json"
	"reflect"
	"regexp"
)

type Request struct {
	Url        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	urlAsRegex regexp.Regexp
}

// Equal checks if a request is identical to another
func (r Request) Equal(r2 Request) bool {
	return reflect.DeepEqual(r, r2)
}

func (r Request) String() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic("cannot marshal request")
	}
	return string(bytes)
}

func WithRequestHeaders(headers map[string]string) func(*Request) {
	return func(r *Request) {
		for k, v := range headers {
			r.Headers[k] = v
		}
	}
}

func WithRequestHeader(header string, value string) func(*Request) {
	return func(r *Request) {
		r.Headers[header] = value
	}
}

func NewRequest(url string, method string, options ...func(*Request)) Request {
	r := Request{
		Url:        url,
		urlAsRegex: *regexp.MustCompile(url),
		Method:     method,
		Headers:    make(map[string]string),
	}
	for _, o := range options {
		o(&r)
	}

	return r
}
