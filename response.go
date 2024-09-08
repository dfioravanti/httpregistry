package servermock

import "net/http"

type Response struct {
	body    []byte
	status  int
	headers map[string]string
}

var Response200 = Response{status: http.StatusOK}

func WithResponseBody(body []byte) func(*Response) {
	return func(r *Response) {
		r.body = body
	}
}

func WithResponseStatus(status int) func(*Response) {
	return func(r *Response) {
		r.status = status
	}
}

func WithResponseHeaders(headers map[string]string) func(*Response) {
	return func(r *Response) {
		for k, v := range headers {
			r.headers[k] = v
		}
	}
}

func WithResponseHeader(header string, value string) func(*Response) {
	return func(r *Response) {
		r.headers[header] = value
	}
}

func NewResponse(options ...func(*Response)) Response {
	r := Response{}
	for _, o := range options {
		o(&r)
	}

	return r
}
