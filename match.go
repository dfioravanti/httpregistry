package httpregistry

import (
	"errors"
	"net/http"
)

var (
	ErrNoNextResponseFound = errors.New("it was not possible to found a next response")
)

// A Match is used to connect a Request to one or multiple possible Response(s) so that when the request happens the mock server
// returns the desired response.
type Match interface {
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

// A FixedResponseMatch is a match that returns a fixed Response each time a predefined Request happens
type FixedResponseMatch struct {
	request       Request
	response      Response
	numberOfCalls int
	matches       []*http.Request
}

// NewFixedResponseMatch creates a new FixedResponseMatch
func NewFixedResponseMatch(request Request, response Response) *FixedResponseMatch {
	return &FixedResponseMatch{
		request:       request,
		response:      response,
		numberOfCalls: 0,
		matches:       make([]*http.Request, 0),
	}
}

// Request returns the request that triggers the match
func (m *FixedResponseMatch) Request() Request {
	return m.request
}

// Next response returns the next response associated with the match and records which request triggered the match.
// It never raises an error
func (m *FixedResponseMatch) NextResponse(req *http.Request) (Response, error) {
	m.matches = append(m.matches, cloneHttpRequest(req))

	return m.response, nil
}

// Matches returns the list of http.Request that matched with this Match
func (m *FixedResponseMatch) Matches() []*http.Request {
	return m.matches
}

// NumberOfCalls returns the number of times the match was fulfilled
func (m *FixedResponseMatch) NumberOfCalls() int {
	return len(m.matches)
}

// A FixedResponseMatch is a match that returns a different Response each time a predefined Request happens
//
// Important: the list of responses gets consumed by the server. Do not reuse this structure, create a new one
type MultipleResponsesMatch struct {
	request       Request
	responses     Responses
	numberOfCalls int
	matches       []*http.Request
}

// NewFixedResponseMatch creates a new FixedResponseMatch
func NewMultipleResponsesMatch(request Request, responses Responses) *MultipleResponsesMatch {
	return &MultipleResponsesMatch{
		request:       request,
		responses:     responses,
		numberOfCalls: 0,
		matches:       make([]*http.Request, 0),
	}
}

// Request returns the request that triggers the match
func (m *MultipleResponsesMatch) Request() Request {
	return m.request
}

// Next response returns the next response associated with the match.
// If the list of responses is exhausted it will return a ErrNoNextResponseFound error
// It consumes the list associated with the MultipleResponsesMatch
func (m *MultipleResponsesMatch) NextResponse(req *http.Request) (Response, error) {
	if len(m.responses) == 0 {
		return Response{}, ErrNoNextResponseFound
	}

	m.matches = append(m.matches, cloneHttpRequest(req))

	head, tail := m.responses[0], m.responses[1:]
	m.responses = tail

	return head, nil
}

// Matches returns the list of http.Request that matched with this Match
func (m *MultipleResponsesMatch) Matches() []*http.Request {
	return m.matches
}

// NumberOfCalls returns the number of times the match was fulfilled
func (m *MultipleResponsesMatch) NumberOfCalls() int {
	return len(m.matches)
}
