package servermock

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Request struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Body       string
	responses  []Response
	urlAsRegex regexp.Regexp
}

func (r Request) String() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic(fmt.Sprintf("cannot marshal request"))
	}
	return string(bytes)
}

func WithRequestURL(url string) func(*Request) {
	return func(r *Request) {
		r.URL = url
		r.urlAsRegex = *regexp.MustCompile(url)
	}
}

func WithRequestMethod(method string) func(*Request) {
	return func(r *Request) {
		r.Method = method
	}
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

func NewRequest(options ...func(*Request)) Request {
	r := Request{}
	for _, o := range options {
		o(&r)
	}

	return r
}
