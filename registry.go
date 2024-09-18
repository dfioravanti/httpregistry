package servermock

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

// A Registry contains all the Match that were registered and it is designed to allow easy access and manipulation to them
type Registry struct {
	matches []Match
	misses  []Miss
}

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	reg := Registry{}
	return &reg
}

// AddSimpleRequest is a helper function for the most common case of wanting to return a 200 response when URL is called with a method
func (reg *Registry) AddSimpleRequest(URL string, method string) {
	request := NewRequest(URL, method)

	reg.matches = append(reg.matches, NewFixedResponseMatch(request, OkResponse))
}

// AddSimpleRequest is a helper function for the common case of wanting to return a statusCode response when URL is called with a method
func (m *Registry) AddSimpleRequestWithStatusCode(URL string, method string, statusCode int) {
	request := NewRequest(URL, method)
	response := NewResponse(statusCode, nil)

	m.matches = append(m.matches, NewFixedResponseMatch(request, response))
}

// AddRequest adds request to the mock server and it returns a 200 response each time that request happens
func (reg *Registry) AddRequest(request Request) {
	reg.matches = append(reg.matches, NewFixedResponseMatch(request, OkResponse))
}

// AddRequest adds request to the mock server and it returns response each time that request happens
func (reg *Registry) AddRequestWithResponse(request Request, response Response) {
	reg.matches = append(reg.matches, NewFixedResponseMatch(request, response))
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

// GetServer returns a httptest.Server designed to match all the requests registered with the Registry
func (reg *Registry) GetServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathToMatch := r.URL.String()
		methodToMatch := r.Method

		for _, possibleMatch := range reg.matches {
			requestToMatch := possibleMatch.Request()
			if requestToMatch.urlAsRegex.MatchString(pathToMatch) {
				response, err := possibleMatch.NextResponse(r)
				if err != nil {
					if errors.Is(ErrNoNextResponseFound, err) {
						t.Errorf("run out of responses when calling: %v %v", requestToMatch.Method, requestToMatch.Url)
					}
				}

				if methodToMatch == r.Method {
					for k, v := range response.headers {
						w.Header().Add(k, v)
					}
					w.WriteHeader(response.status)
					_, err = w.Write(response.body)
					if err != nil {
						panic("cannot write body of request")
					}

					return
				} else {
					miss := Miss{
						MissedMatch: requestToMatch,
						Why:         methodDoesNotMatch,
					}
					reg.misses = append(reg.misses, miss)
				}
			} else {
				miss := Miss{
					MissedMatch: requestToMatch,
					Why:         pathDoesNotMatch,
				}
				reg.misses = append(reg.misses, miss)
			}
		}

		res, err := httputil.DumpRequest(r, true)
		if err != nil {
			t.Errorf("impossible to dump http request with error %v", err)
		}

		t.Errorf("no registered request matched %v, you can use .Why() to get an explanation of why", res)
	}))
}
