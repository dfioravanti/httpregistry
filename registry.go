package httpregistry

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

// A Registry contains all the Match that were registered and it is designed to allow easy access and manipulation to them
type Registry struct {
	matches []match
	misses  []miss
}

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	reg := Registry{}
	return &reg
}

// AddSimpleRequest is a helper function for the most common case of wanting to return a 200 response when URL is called with a method
func (reg *Registry) AddSimpleRequest(URL string, method string) {
	request := NewRequest(URL, method)

	reg.matches = append(reg.matches, newFixedResponseMatch(request, OkResponse))
}

// AddSimpleRequest is a helper function for the common case of wanting to return a statusCode response when URL is called with a method
func (reg *Registry) AddSimpleRequestWithStatusCode(URL string, method string, statusCode int) {
	request := NewRequest(URL, method)
	response := NewResponse(statusCode, nil)

	reg.matches = append(reg.matches, newFixedResponseMatch(request, response))
}

// AddRequest adds request to the mock server and it returns a 200 response each time that request happens
func (reg *Registry) AddRequest(request Request) {
	reg.matches = append(reg.matches, newFixedResponseMatch(request, OkResponse))
}

// AddRequest adds request to the mock server and it returns response each time that request happens
func (reg *Registry) AddRequestWithResponse(request Request, response Response) {
	reg.matches = append(reg.matches, newFixedResponseMatch(request, response))
}

func (reg *Registry) AddRequestWithResponses(request Request, responses ...Response) {
	reg.matches = append(reg.matches, newMultipleResponsesMatch(request, responses))
}

func (reg *Registry) GetMatchesPerRequest(r Request) []*http.Request {
	for _, match := range reg.matches {
		if match.Request().Equal(r) {
			return match.Matches()
		}
	}

	return []*http.Request{}
}

func (reg *Registry) GetMatchesUrl(url string) []*http.Request {
	for _, match := range reg.matches {
		r := match.Request()
		if r.urlAsRegex.MatchString(url) {
			return match.Matches()
		}
	}

	return []*http.Request{}
}

func (reg *Registry) GetMatchesUrlAndMethod(url string, method string) []*http.Request {
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
	expectedUrlAsRegex := registeredMatch.Request().urlAsRegex
	expectedMethod := registeredMatch.Request().Method
	expectedHeaders := registeredMatch.Request().Headers

	if !expectedUrlAsRegex.MatchString(r.URL.String()) {
		return false, []miss{{
			MissedMatch: registeredMatch,
			Why:         pathDoesNotMatch,
		}}

	}

	if expectedMethod != r.Method {
		return false, []miss{{
			MissedMatch: registeredMatch,
			Why:         methodDoesNotMatch,
		}}
	}

	headersMisses := []miss{}
	for headerToMatch, valueToMatch := range expectedHeaders {
		value := r.Header.Get(headerToMatch)
		if value == "" || value != valueToMatch {
			miss := miss{
				MissedMatch: registeredMatch,
				Why:         headersDoNotMatch,
			}
			headersMisses = append(headersMisses, miss)
		}
	}

	if len(headersMisses) != 0 {
		return false, headersMisses
	}

	return true, nil
}

// GetServer returns a httptest.Server designed to match all the requests registered with the Registry
func (reg *Registry) GetServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We reset the misses since if a previous request matched it is pointless to record that some of the mocks did not match it.
		// If said request did not match then the test would have crashed in any case so the information in misses is useless.
		reg.misses = []miss{}
		for _, possibleMatch := range reg.matches {
			requestToMatch := possibleMatch.Request()

			doesMatch, misses := doesRegisteredMatchMatchIncomingRequest(possibleMatch, r)
			if !doesMatch {
				reg.misses = append(reg.misses, misses...)
				continue
			}

			response, err := possibleMatch.NextResponse(r)
			if err != nil {
				if errors.Is(errNoNextResponseFound, err) {
					t.Errorf("run out of responses when calling: %v %v", requestToMatch.Method, requestToMatch.Url)
				}
			}

			for k, v := range response.headers {
				w.Header().Add(k, v)
			}
			w.WriteHeader(response.status)
			_, err = w.Write(response.body)
			if err != nil {
				panic("cannot write body of request")
			}
			return
		}

		res, err := httputil.DumpRequest(r, true)
		if err != nil {
			t.Errorf("impossible to dump http request with error %v", err)
		}

		t.Errorf("no registered request matched %v\n you can use .Why() to get an explanation of why", string(res))
	}))
}

// Why returns a string that contains all the reasons why the request submitted to the registry failed to match with the registered requests.
// The envision use of this function is just as a helper when debugging the tests,
// most of the time it might not be obvious if there is a typo or a small error.
func (reg *Registry) Why() string {
	output := ""
	for _, miss := range reg.misses {
		output += miss.String() + "\n"
	}
	return output
}
