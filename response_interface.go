package httpregistry

import "net/http"

// mockResponse represents all structs inside the httpregistry that can be used to define how a mocked response will look like.
// Currently we have
//
//   - httpregistry.Response -> it allows to define (one or more) status code, body and headers
//   - httpregistry.CustomResponse -> it allows to define the response as a function of (w, r)
type mockResponse interface {
	// createResponse emits the response encoded in the struct that implements mockResponse to w
	createResponse(w http.ResponseWriter, r *http.Request)
	// String Marshals the mockResponse into string
	String() string
}

type mockResponses = []mockResponse
