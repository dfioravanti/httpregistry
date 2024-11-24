package httpregistry

import (
	"errors"
	"net/http"
)

var (
	errNoNextResponseFound = errors.New("it was not possible to found a next response")
)

// A match is used to connect a Request to one or multiple possible Response(s) so that when the request happens the mock server
// returns the desired response.
type match interface {
	// Request returns the request that triggers the match
	Request() Request
	// Next response returns the next response associated with the match and records which request triggered the match.
	// If the list of responses is exhausted it will return a ErrNoNextResponseFound error
	NextResponse(req *http.Request) (Response, error)
	// NumberOfCalls returns the number of times the match was fulfilled
	NumberOfCalls() int
	// Matches returns the list of http.Request that matched with this Match
	Matches() []*http.Request
}

// A fixedResponseMatch is a match that returns a fixed Response each time a predefined Request happens
type fixedResponseMatch struct {
	request       Request
	response      Response
	numberOfCalls int
	matches       []*http.Request
}

// newFixedResponseMatch creates a new FixedResponseMatch
func newFixedResponseMatch(request Request, response Response) *fixedResponseMatch {
	return &fixedResponseMatch{
		request:       request,
		response:      response,
		numberOfCalls: 0,
		matches:       make([]*http.Request, 0),
	}
}

// Request returns the request that triggers the match
func (m *fixedResponseMatch) Request() Request {
	return m.request
}

// Next response returns the next response associated with the match and records which request triggered the match.
// It never raises an error
func (m *fixedResponseMatch) NextResponse(req *http.Request) (Response, error) {
	m.matches = append(m.matches, cloneHTTPRequest(req))

	return m.response, nil
}

// Matches returns the list of http.Request that matched with this Match
func (m *fixedResponseMatch) Matches() []*http.Request {
	return m.matches
}

// NumberOfCalls returns the number of times the match was fulfilled
func (m *fixedResponseMatch) NumberOfCalls() int {
	return len(m.matches)
}

// A multipleResponsesMatch is a match that returns a different Response each time a predefined Request happens
//
// Important: the list of responses gets consumed by the server. Do not reuse this structure, create a new one
type multipleResponsesMatch struct {
	request       Request
	responses     Responses
	numberOfCalls int
	matches       []*http.Request
}

// NewFixedResponseMatch creates a new FixedResponseMatch
func newMultipleResponsesMatch(request Request, responses Responses) *multipleResponsesMatch {
	return &multipleResponsesMatch{
		request:       request,
		responses:     responses,
		numberOfCalls: 0,
		matches:       make([]*http.Request, 0),
	}
}

// Request returns the request that triggers the match
func (m *multipleResponsesMatch) Request() Request {
	return m.request
}

// Next response returns the next response associated with the match.
// If the list of responses is exhausted it will return a ErrNoNextResponseFound error
// It consumes the list associated with the MultipleResponsesMatch
func (m *multipleResponsesMatch) NextResponse(req *http.Request) (Response, error) {
	if len(m.responses) == 0 {
		return Response{}, errNoNextResponseFound
	}

	m.matches = append(m.matches, cloneHTTPRequest(req))

	head, tail := m.responses[0], m.responses[1:]
	m.responses = tail

	return head, nil
}

// Matches returns the list of http.Request that matched with this Match
func (m *multipleResponsesMatch) Matches() []*http.Request {
	return m.matches
}

// NumberOfCalls returns the number of times the match was fulfilled
func (m *multipleResponsesMatch) NumberOfCalls() int {
	return len(m.matches)
}
