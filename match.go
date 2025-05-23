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
	NextResponse() (mockResponse, error)
	// NumberOfCalls returns the number of times the match was fulfilled
	NumberOfCalls() int
	// Matches returns the list of http.Request that matched with this Match
	Matches() []*http.Request
}

// A consumableResponsesMatch is a match that returns a different Response each time a predefined Request happens
// Important: the list of responses gets consumed by the server. Do not reuse this structure, create a new one
type consumableResponsesMatch struct {
	request       Request
	responses     mockResponses
	numberOfCalls int
	matches       []*http.Request
}

// newConsumableResponsesMatch creates a new consumableResponsesMatch
func newConsumableResponsesMatch(request Request, responses mockResponses) *consumableResponsesMatch {
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
func (m *consumableResponsesMatch) NextResponse() (mockResponse, error) {
	if len(m.responses) == 0 {
		return nil, errNoNextResponseFound
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

// A infiniteResponsesMatch is a match that returns the same Response each time a predefined Request happens.
// The response is never consumed, so NextResponse() never returns an errNoNextResponseFound
type infiniteResponsesMatch struct {
	request       Request
	response      mockResponse
	numberOfCalls int
	matches       []*http.Request
}

// newInfiniteResponsesMatch creates a new infiniteResponsesMatch
func newInfiniteResponsesMatch(request Request, response mockResponse) *infiniteResponsesMatch {
	return &infiniteResponsesMatch{
		request:       request,
		response:      response,
		numberOfCalls: 0,
		matches:       make([]*http.Request, 0),
	}
}

// Request returns the request that triggers the match
func (m *infiniteResponsesMatch) Request() Request {
	return m.request
}

// RecordMatch records that a request was a successful match for this match
func (m *infiniteResponsesMatch) RecordMatch(req *http.Request) {
	m.matches = append(m.matches, cloneHTTPRequest(req))
}

// Next response returns the next response associated with the match.
// As this is an infinite match this function never returns an error
func (m *infiniteResponsesMatch) NextResponse() (mockResponse, error) {
	return m.response, nil
}

// Matches returns the list of http.Request that matched with this Match
func (m *infiniteResponsesMatch) Matches() []*http.Request {
	return m.matches
}

// NumberOfCalls returns the number of times the match was fulfilled
func (m *infiniteResponsesMatch) NumberOfCalls() int {
	return len(m.matches)
}
