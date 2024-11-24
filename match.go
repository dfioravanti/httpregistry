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
	// RecordMatch records that a request was a successful match for this match
	RecordMatch(req *http.Request)
	// Next response returns the next response associated with the match and records which request triggered the match.
	// If the list of responses is exhausted it will return a ErrNoNextResponseFound error
	NextResponse() (Response, error)
	// NumberOfCalls returns the number of times the match was fulfilled
	NumberOfCalls() int
	// Matches returns the list of http.Request that matched with this Match
	Matches() []*http.Request
}

// A consumableResponsesMatch is a match that returns a different Response each time a predefined Request happens
// Important: the list of responses gets consumed by the server. Do not reuse this structure, create a new one
type consumableResponsesMatch struct {
	request       Request
	responses     Responses
	numberOfCalls int
	matches       []*http.Request
}

// newConsumableResponsesMatch creates a new consumableResponsesMatch
func newConsumableResponsesMatch(request Request, responses Responses) *consumableResponsesMatch {
	return &consumableResponsesMatch{
		request:       request,
		responses:     responses,
		numberOfCalls: 0,
		matches:       make([]*http.Request, 0),
	}
}

// Request returns the request that triggers the match
func (m *consumableResponsesMatch) Request() Request {
	return m.request
}

// RecordMatch records that a request was a successful match for this match
func (m *consumableResponsesMatch) RecordMatch(req *http.Request) {
	m.matches = append(m.matches, cloneHTTPRequest(req))
}

// Next response returns the next response associated with the match.
// If the list of responses is exhausted it will return a ErrNoNextResponseFound error
// It consumes the list associated with the MultipleResponsesMatch
func (m *consumableResponsesMatch) NextResponse() (Response, error) {
	if len(m.responses) == 0 {
		return Response{}, errNoNextResponseFound
	}

	head, tail := m.responses[0], m.responses[1:]
	m.responses = tail

	return head, nil
}

// Matches returns the list of http.Request that matched with this Match
func (m *consumableResponsesMatch) Matches() []*http.Request {
	return m.matches
}

// NumberOfCalls returns the number of times the match was fulfilled
func (m *consumableResponsesMatch) NumberOfCalls() int {
	return len(m.matches)
}
