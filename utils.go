package httpregistry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// cloneHTTPRequest clones a http.Request in full.
// By default the .clone does not clone the body
func cloneHTTPRequest(req *http.Request) *http.Request {
	var buf []byte
	var err error
	if req.Body != nil {
		buf, err = io.ReadAll(req.Body)
		if err != nil {
			panic(fmt.Sprintf("cannot read body of request with error: %v", err))
		}
	}

	newRequest := req.Clone(req.Context())

	req.Body = io.NopCloser(bytes.NewBuffer(buf))
	newRequest.Body = io.NopCloser(bytes.NewBuffer(buf))

	return newRequest
}
