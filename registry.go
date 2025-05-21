package httpregistry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"reflect"
)

// Registry represents a collection of matches that associate to a http request a http response.
// It contains all the Match that were registered and after the server is called it contains all the reasons why a request did not match with a particular match
// the testing.T is used to signal that there was an unexpected error or that not all the responses were consumed as expected
type Registry struct {
	t       TestingT
	matches []match
	misses  []miss
}

// NewRegistry creates a new empty Registry
func NewRegistry(t TestingT) *Registry {
	reg := Registry{t: t}
	return &reg
}

// Add adds to the registry a 200 response for any requests
//
//	reg := httpregistry.NewRegistry(t)
//	reg.Add()
//	reg.GetServer()
//
// will create a http server that returns 200 on calling anything.
func (reg *Registry) Add() {
	request := NewRequest()
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{OkResponse}))
}

// AddURL adds to the registry a 200 response for a request that matches the URL
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddURL("/foo")
//	reg.GetServer()
//
// will create a http server that returns 200 on calling GET "/foo" and fails the test on anything else
func (reg *Registry) AddURL(URL string) {
	request := NewRequest().WithURL(URL)
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{OkResponse}))
}

// AddURLWithStatusCode adds to the registry a statusCode response for a request that matches the URL
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddURLWithStatusCode("/foo", 401)
//	reg.GetServer()
//
// will create a http server that returns 401 on calling GET "/foo" and fails the test on anything else
func (reg *Registry) AddURLWithStatusCode(URL string, statusCode int) {
	request := NewRequest().WithURL(URL)
	response := NewResponse().WithStatus(statusCode)

	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{response}))
}

// AddMethod adds to the registry a 200 response for a request that matches the method
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddMethod("/foo")
//	reg.GetServer()
//
// will create a http server that returns 200 on calling GET "/foo" and fails the test on anything else
func (reg *Registry) AddMethod(method string) {
	request := NewRequest().WithMethod(method)
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{OkResponse}))
}

// AddMethodWithStatusCode adds to the registry a statusCode response for a request that matches the method
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddMethodWithStatusCode("/foo", 401)
//	reg.GetServer()
//
// will create a http server that returns 401 on calling GET "/foo" and fails the test on anything else
func (reg *Registry) AddMethodWithStatusCode(method string, statusCode int) {
	request := NewRequest().WithMethod(method)
	response := NewResponse().WithStatus(statusCode)
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{response}))
}

// AddMethodAndURL adds to the registry a 200 response for a request that matches method and URL
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddMethodAndURL(GET, "/foo")
//	reg.GetServer()
//
// will create a http server that returns 200 on calling GET "/foo" and fails the test on anything else
func (reg *Registry) AddMethodAndURL(method string, URL string) {
	request := NewRequest().WithMethod(method).WithURL(URL)
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{OkResponse}))
}

// AddMethodAndURLWithStatusCode adds to the registry a statusCode response for a request that matches method and URL
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddSimpleRequest(PUT, "/foo", 204)
//	reg.GetServer()
//
// will create a http server that returns 204 on calling GET "/foo" and fails the test on anything else
func (reg *Registry) AddMethodAndURLWithStatusCode(method string, URL string, statusCode int) {
	request := NewRequest().WithMethod(method).WithURL(URL)
	response := NewResponse().WithStatus(statusCode)

	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{response}))
}

// AddBody adds to the registry a statusCode response for a request that matches method and URL
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddSimpleRequest(PUT, "/foo", 204)
//	reg.GetServer()
//
// will create a http server that returns 204 on calling GET "/foo" and fails the test on anything else
func (reg *Registry) AddBody(body []byte) {
	request := NewRequest().WithBody(body)
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{OkResponse}))
}

// AddRequest adds to the registry a 200 response for a generic request that needs to be matched
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddRequest(
//		httpregistry.NewRequest(GET, "/foo", httpregistry.WithRequestHeader("header", "value"))
//	)
//	reg.GetServer()
//
// will create a http server that returns 200 on calling GET "/foo" with the correct header and fails the test on anything else
func (reg *Registry) AddRequest(request Request) {
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{OkResponse}))
}

// AddResponse adds to the registry a generic response that is returned for any call
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddResponse(
//		httpregistry.NewResponse(http.StatusCreated, []byte{"hello"}),
//	)
//	reg.GetServer()
//
// will create a http server that returns 204 with "hello" as body on calling the server on any URL
func (reg *Registry) AddResponse(response mockResponse) {
	request := NewRequest()
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{response}))
}

// AddResponses adds to the registry a generic response that is returned for any call
//
//		reg := httpregistry.NewRegistry(t)
//		reg.AddResponses(
//			httpregistry.NewResponse(http.StatusCreated, []byte{"hello"}),
//	     httpregistry.NewResponse(http.StatusCreated, []byte{"hello"}),
//		)
//		reg.GetServer()
//
// will create a http server that returns 204 with "hello" as body on calling the server on any URL for two times and then returns an error
func (reg *Registry) AddResponses(responses ...mockResponse) {
	request := NewRequest()
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, responses))
}

// AddRequestWithResponse adds to the registry a generic response for a generic request that needs to be matched
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddRequest(
//		httpregistry.NewRequest(GET, "/foo", httpregistry.WithRequestHeader("header", "value")),
//		httpregistry.NewResponse(http.StatusCreated, []byte{"hello"}),
//	)
//	reg.GetServer()
//
// will create a http server that returns 204 with "hello" as body on calling GET "/foo" with the correct header and fails the test on anything else
func (reg *Registry) AddRequestWithResponse(request Request, response mockResponse) {
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, mockResponses{response}))
}

// AddRequestWithResponses adds to the registry multiple responses for a generic request that needs to be matched.
// The responses are consumed by the calls so if more calls than responses will happen then the test will fail
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddRequestWithResponses(
//		httpregistry.NewRequest(GET, "/foo", httpregistry.WithRequestHeader("header", "value")),
//		httpregistry.NewResponse(http.StatusCreated, []byte{"hello"}),
//		httpregistry.NewResponse(http.Ok, []byte{"hello again"}),
//	)
//	reg.GetServer()
//
// will create a http server that returns 204 with "hello" as body on calling GET "/foo" the first call with the correct header,
// it returns 200 with "hello again" as body on the second call with the correct header and fails the test on anything else
func (reg *Registry) AddRequestWithResponses(request Request, responses ...mockResponse) {
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, responses))
}

// GetMatchesForRequest returns the *http.Request that matched a generic Request
func (reg *Registry) GetMatchesForRequest(r Request) []*http.Request {
	for _, match := range reg.matches {
		if match.Request().Equal(r) {
			return match.Matches()
		}
	}
	return []*http.Request{}
}

// GetMatchesForURL returns the http.Requests that matched a specific URL independently of the method used to call it
func (reg *Registry) GetMatchesForURL(url string) []*http.Request {
	for _, match := range reg.matches {
		r := match.Request()
		if r.urlAsRegex.MatchString(url) {
			return match.Matches()
		}
	}
	return []*http.Request{}
}

// GetMatchesURLAndMethod returns the http.Requests that matched a specific method, URL pair
func (reg *Registry) GetMatchesURLAndMethod(url string, method string) []*http.Request {
	for _, match := range reg.matches {
		r := match.Request()
		if r.urlAsRegex.MatchString(url) && r.Method == method {
			return match.Matches()
		}
	}
	return []*http.Request{}
}

// doesRegisteredMatchMatchIncomingRequest checks if the incoming request is a match for the match that we are currently evaluating.
// If it is not a match this function will return a slice of miss objects that explain why the match is not possible.
func doesRegisteredMatchMatchIncomingRequest(registeredMatch match, r *http.Request) (bool, []miss) {
	// The default request matches everything so no point in checking further
	if registeredMatch.Request().Equal(NewRequest()) {
		return true, nil
	}

	expectedURL := registeredMatch.Request().URL
	expectedURLAsRegex := registeredMatch.Request().urlAsRegex
	expectedMethod := registeredMatch.Request().Method
	expectedHeaders := registeredMatch.Request().Headers
	expectedBody := registeredMatch.Request().Body

	misses := []miss{}

	// if the match contains the default values then there is no point in saying that something was missed
	if expectedURL != "" {
		if expectedURLAsRegex.MatchString(r.URL.String()) {
			return true, nil
		}
		misses = append(misses, newMiss(registeredMatch, pathDoesNotMatch))
	}

	if expectedMethod != "" {
		if expectedMethod == r.Method {
			return true, nil
		}
		misses = append(misses, newMiss(registeredMatch, methodDoesNotMatch))
	}

	if len(expectedBody) > 0 {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(fmt.Errorf("cannot read the body of the request: %w", err))
		}
		if reflect.DeepEqual(expectedBody, body) {
			return true, nil
		}
	}

	if len(expectedHeaders) > 0 {
		headersMisses := []miss{}
		for headerToMatch, valueToMatch := range expectedHeaders {
			value := r.Header.Get(headerToMatch)
			if value == "" || value != valueToMatch {
				misses = append(misses, newMiss(registeredMatch, headerDoesNotMatch))
			}
		}

		if len(headersMisses) == 0 {
			return true, nil
		}
		misses = append(misses, headersMisses...)
	}

	return false, misses
}

// GetServer returns a httptest.Server designed to match all the requests registered with the Registry
func (reg *Registry) GetServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We reset the misses since if a previous request matched it is pointless to record that some of the mocks did not match it.
		// If said request did not match then the test would have crashed in any case so the information in misses is useless.
		reg.misses = []miss{}
		for _, possibleMatch := range reg.matches {
			doesMatch, misses := doesRegisteredMatchMatchIncomingRequest(possibleMatch, r)
			if !doesMatch {
				reg.misses = append(reg.misses, misses...)
				continue
			}

			response, err := possibleMatch.NextResponse()
			if err != nil {
				if errors.Is(errNoNextResponseFound, err) {
					reg.misses = append(reg.misses, newMiss(possibleMatch, outOfResponses))
					continue
				}
			}

			possibleMatch.RecordMatch(r)
			response.createResponse(w, r)
			return
		}

		res, err := httputil.DumpRequest(r, true)
		if err != nil {
			reg.t.Errorf("impossible to dump http request with error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		reg.t.Errorf("no registered request matched %v\n The reasons why this is the case are returned in the body", string(res))
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(reg.Why()))
	}))
}

// CheckAllResponsesAreConsumed fails the test if there are unused responses at the end of the test.
// This is useful to check if all the expected calls happened or if there is an unexpected behavior happening.
func (reg *Registry) CheckAllResponsesAreConsumed() {
	for _, match := range reg.matches {
		response, err := match.NextResponse()
		if err == nil {
			reg.t.Errorf("request %v has %v as unused response", match.Request(), response)
		}
	}
}

// Why returns a string that contains all the reasons why the request submitted to the registry failed to match with the registered requests.
// The envision use of this function is just as a helper when debugging the tests,
// most of the time it might not be obvious if there is a typo or a small error.
func (reg *Registry) Why() string {
	output, err := json.Marshal(reg.misses)
	if err != nil {
		reg.t.Errorf("impossible to serialize matches to json: %v", err)
	}
	return string(output)
}
