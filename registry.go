package httpregistry

import (
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
	t                          TestingT
	matches                    []match
	misses                     []miss
	nameRequestFunction        func() string
	nameCustomResponseFunction func() string
	nameResponseFunction       func() string
}

// NewRegistry creates a new empty Registry
func NewRegistry(t TestingT) *Registry {
	reg := Registry{
		t:                          t,
		nameRequestFunction:        defaultName("mock request"),
		nameCustomResponseFunction: defaultName("custom mock response"),
		nameResponseFunction:       defaultName("mock response"),
	}
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
	request := NewRequest().WithName(reg.nameRequestFunction())
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
	request := NewRequest().WithURL(URL).WithName(reg.nameRequestFunction())
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
	request := NewRequest().WithURL(URL).WithName(reg.nameRequestFunction())
	response := NewResponse().WithStatus(statusCode).WithName(reg.nameResponseFunction())

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
	request := NewRequest().WithMethod(method).WithName(reg.nameRequestFunction())
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
	request := NewRequest().WithMethod(method).WithName(reg.nameRequestFunction())
	response := NewResponse().WithStatus(statusCode).WithName(reg.nameResponseFunction())
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
	request := NewRequest().WithMethod(method).WithURL(URL).WithName(reg.nameRequestFunction())
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
	request := NewRequest().WithMethod(method).WithURL(URL).WithName(reg.nameRequestFunction())
	response := NewResponse().WithStatus(statusCode).WithName(reg.nameResponseFunction())

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
	request := NewRequest().WithBody(body).WithName(reg.nameRequestFunction())
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
	request = reg.ifNeededSetDefaultNameToRequest(request)
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
	request := NewRequest().WithName(reg.nameRequestFunction())
	response = reg.ifNeededSetDefaultNameToMockResponse(response)
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
	request := NewRequest().WithName(reg.nameRequestFunction())
	responsesWithNames := make(mockResponses, 0, len(responses))
	for _, response := range responses {
		responsesWithNames = append(responsesWithNames, reg.ifNeededSetDefaultNameToMockResponse(response))
	}
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, responsesWithNames))
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
	request = reg.ifNeededSetDefaultNameToRequest(request)
	response = reg.ifNeededSetDefaultNameToMockResponse(response)
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
	request = reg.ifNeededSetDefaultNameToRequest(request)
	responsesWithNames := make(mockResponses, 0, len(responses))
	for _, response := range responses {
		responsesWithNames = append(responsesWithNames, reg.ifNeededSetDefaultNameToMockResponse(response))
	}
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, responsesWithNames))
}

// AddInfiniteResponse adds to the registry a generic response that is returned for any call and it is never consumed
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddInfiniteResponse(
//		httpregistry.NewResponse(http.StatusCreated, []byte{"hello"}),
//	)
//	reg.GetServer()
//
// will create a http server that returns 204 with "hello" as body on calling the server on any URL for as many times as needed
func (reg *Registry) AddInfiniteResponse(response mockResponse) {
	request := NewRequest().WithName(reg.nameRequestFunction())
	response = reg.ifNeededSetDefaultNameToMockResponse(response)
	reg.matches = append(reg.matches, newInfiniteResponsesMatch(request, response))
}

// AddRequestWithInfiniteResponse adds to the registry a generic response for a generic request that needs to be matched
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddRequestWithInfiniteResponse(
//		httpregistry.NewRequest(GET, "/foo", httpregistry.WithRequestHeader("header", "value")),
//		httpregistry.NewResponse(http.StatusCreated, []byte{"hello"}),
//	)
//	reg.GetServer()
//
// will create a http server that returns 204 with "hello" as body on calling GET "/foo" with the correct header and fails the test on anything else
func (reg *Registry) AddRequestWithInfiniteResponse(request Request, response mockResponse) {
	request = reg.ifNeededSetDefaultNameToRequest(request)
	response = reg.ifNeededSetDefaultNameToMockResponse(response)
	reg.matches = append(reg.matches, newInfiniteResponsesMatch(request, response))
}

// GetMatchesForRequest returns the *http.Request that matched a generic Request
func (reg *Registry) GetMatchesForRequest(r Request) []*http.Request {
	for _, match := range reg.matches {
		if match.Request().Equal(r) {
			matches := match.Matches()

			// we clone the requests so that if this function is called multiple times things
			// like the request body can be accessed again
			matchesToReturn := make([]*http.Request, 0, len(matches))
			for _, m := range matches {
				matchesToReturn = append(matchesToReturn, cloneHTTPRequest(m))
			}
			return matchesToReturn
		}
	}
	return []*http.Request{}
}

// GetMatchesForURL returns the http.Requests that matched a specific URL independently of the method used to call it
func (reg *Registry) GetMatchesForURL(url string) []*http.Request {
	for _, match := range reg.matches {
		r := match.Request()
		if r.urlAsRegex.MatchString(url) {
			matches := match.Matches()

			// we clone the requests so that if this function is called multiple times things
			// like the request body can be accessed again
			matchesToReturn := make([]*http.Request, 0, len(matches))
			for _, m := range matches {
				matchesToReturn = append(matchesToReturn, cloneHTTPRequest(m))
			}
			return matchesToReturn
		}
	}
	return []*http.Request{}
}

// GetMatchesURLAndMethod returns the http.Requests that matched a specific method, URL pair
func (reg *Registry) GetMatchesURLAndMethod(url string, method string) []*http.Request {
	for _, match := range reg.matches {
		r := match.Request()
		if r.urlAsRegex.MatchString(url) && r.method == method {
			matches := match.Matches()

			// we clone the requests so that if this function is called multiple times things
			// like the request body can be accessed again
			matchesToReturn := make([]*http.Request, 0, len(matches))
			for _, m := range matches {
				matchesToReturn = append(matchesToReturn, cloneHTTPRequest(m))
			}
			return matchesToReturn
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

	expectedURL := registeredMatch.Request().url
	expectedURLAsRegex := registeredMatch.Request().urlAsRegex
	expectedMethod := registeredMatch.Request().method
	expectedHeaders := registeredMatch.Request().headers
	expectedBody := registeredMatch.Request().body

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
			response.serveResponse(w, r)
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
//
// **Important**: If you are using AddInfiniteRequest this call will ALWAYS fail!
func (reg *Registry) CheckAllResponsesAreConsumed() {
	for _, match := range reg.matches {
		response, err := match.NextResponse()
		if err == nil {
			reg.t.Errorf("request %v has %v as unused response", match.Request().String(), response)
		}
	}
}

// Why returns a string that contains all the reasons why the request submitted to the registry failed to match with the registered requests.
// The envision use of this function is just as a helper when debugging the tests,
// most of the time it might not be obvious if there is a typo or a small error.
func (reg *Registry) Why() string {
	outputString := ""
	for i, miss := range reg.misses {
		if i == 0 {
			outputString = miss.String()
		} else {
			outputString += "\n" + miss.String()
		}
	}
	return outputString
}

// ifNeededSetDefaultNameToRequest overwrites the name field in a Request if the name is currently the empty string
func (reg Registry) ifNeededSetDefaultNameToRequest(request Request) Request {
	if request.name == "" {
		request = request.WithName(reg.nameRequestFunction())
	}
	return request
}

// ifNeededSetDefaultNameToResponse overwrites the name field in a Response if the name is currently the empty string
func (reg Registry) ifNeededSetDefaultNameToResponse(response Response) Response {
	if response.name == "" {
		response = response.WithName(reg.nameResponseFunction())
	}
	return response
}

// ifNeededSetDefaultNameToCustomRequest overwrites the name field in a CustomResponse if the name is currently the empty string
func (reg Registry) ifNeededSetDefaultNameToCustomResponse(response CustomResponse) CustomResponse {
	if response.name == "" {
		response = response.WithName(reg.nameCustomResponseFunction())
	}
	return response
}

// ifNeededSetDefaultNameToMockResponse overwrites the name field in a mockResponse if the name is currently the empty string
func (reg Registry) ifNeededSetDefaultNameToMockResponse(response mockResponse) mockResponse {
	switch r := response.(type) {
	case Response:
		response = reg.ifNeededSetDefaultNameToResponse(r)
	case CustomResponse:
		response = reg.ifNeededSetDefaultNameToCustomResponse(r)
	}

	return response
}
