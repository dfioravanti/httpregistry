package httpregistry

import (
	"bytes"
	"encoding/json"
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

// mustMarshalJSON tries to marshal v into JSON and panics if it cannot
func mustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("body cannot be marshaled to JSON: %s", err))
	}
	return b
}

// defaultName is used create default names for requests and responses
func defaultName(baseString string) func() string {
	counter := 1
	return func() string {
		stringToReturn := fmt.Sprintf("%s #%d", baseString, counter)
		counter++
		return stringToReturn
	}
}
