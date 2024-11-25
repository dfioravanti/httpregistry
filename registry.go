package httpregistry

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
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

// AddSimpleRequest adds to the registry a 200 response for a request that matches method and URL
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddSimpleRequest(GET, "/foo")
//	reg.GetServer()
//
// will create a http server that returns 200 on calling GET "/foo" and 404 on anything else
func (reg *Registry) AddSimpleRequest(method string, URL string) {
	request := NewRequest(method, URL)
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, []Response{OkResponse}))
}

// AddSimpleRequestWithStatusCode adds to the registry a statusCode response for a request that matches method and URL
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddSimpleRequest(PUT, "/foo", 204)
//	reg.GetServer()
//
// will create a http server that returns 204 on calling GET "/foo" and 404 on anything else
func (reg *Registry) AddSimpleRequestWithStatusCode(method string, URL string, statusCode int) {
	request := NewRequest(method, URL)
	response := NewResponse(statusCode, nil)

	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, []Response{response}))
}

// AddRequest adds to the registry a 200 response for a generic request that needs to be matched
//
//	reg := httpregistry.NewRegistry(t)
//	reg.AddRequest(
//		httpregistry.NewRequest(GET, "/foo", httpregistry.WithRequestHeader("header", "value"))
//	)
//	reg.GetServer()
//
// will create a http server that returns 200 on calling GET "/foo" with the correct header and 404 on anything else
func (reg *Registry) AddRequest(request Request) {
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, []Response{OkResponse}))
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
// will create a http server that returns 204 with "hello" as body on calling GET "/foo" with the correct header and 404 on anything else
func (reg *Registry) AddRequestWithResponse(request Request, response Response) {
	reg.matches = append(reg.matches, newConsumableResponsesMatch(request, []Response{response}))
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
// it returns 200 with "hello again" as body on the second call with the correct header and 404 on anything else
func (reg *Registry) AddRequestWithResponses(request Request, responses ...Response) {
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

// GetMatchesForURL returns the *http.Request that matched a specific URL independently of the method used to call it
func (reg *Registry) GetMatchesForURL(url string) []*http.Request {
	for _, match := range reg.matches {
		r := match.Request()
		if r.urlAsRegex.MatchString(url) {
			return match.Matches()
		}
	}
	return []*http.Request{}
}

// GetMatchesURLAndMethod returns the *http.Request that matched a specific method, URL pair
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
	expectedURLAsRegex := registeredMatch.Request().urlAsRegex
	expectedMethod := registeredMatch.Request().Method
	expectedHeaders := registeredMatch.Request().Headers

	if !expectedURLAsRegex.MatchString(r.URL.String()) {
		return false, []miss{newMiss(registeredMatch, pathDoesNotMatch)}
	}

	if expectedMethod != r.Method {
		return false, []miss{newMiss(registeredMatch, methodDoesNotMatch)}
	}

	headersMisses := []miss{}
	for headerToMatch, valueToMatch := range expectedHeaders {
		value := r.Header.Get(headerToMatch)
		if value == "" || value != valueToMatch {
			headersMisses = append(headersMisses, newMiss(registeredMatch, headerDoesNotMatch))
		}
	}

	if len(headersMisses) != 0 {
		return false, headersMisses
	}

	return true, nil
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

			for k, v := range response.Headers {
				w.Header().Add(k, v)
			}
			w.WriteHeader(response.Status)
			_, err = w.Write(response.Body)
			if err != nil {
				panic("cannot write body of request")
			}
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
