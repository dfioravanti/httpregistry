package httpregistry

import (
	"net/http"
)

func (s *TestSuite) TestFixedResponseMatchHasExpectedResponse() {
	request := NewRequest("/", http.MethodPost)
	response := NoContentResponse

	match := newFixedResponseMatch(request, response)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestFixedResponseRespondsForever() {
	request := NewRequest("/", http.MethodPost)
	expectedResponse := NoContentResponse

	match := newFixedResponseMatch(request, expectedResponse)
	expectedNumberOfCalls := 1000

	for range expectedNumberOfCalls {
		response, err := match.NextResponse(&http.Request{})

		s.NoError(err)
		s.Equal(request, match.Request())
		s.Equal(expectedResponse, response)
	}

	s.Equal(expectedNumberOfCalls, match.NumberOfCalls())
}

func (s *TestSuite) TestMultipleResponsesHasExpectedResponse() {
	request := NewRequest("/", http.MethodPost)
	responses := Responses{NoContentResponse}

	match := newMultipleResponsesMatch(request, responses)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestMultipleResponsesRespondsTheCorrectNumberOfTimes() {
	request := NewRequest("/", http.MethodPost)

	expectedFirstResponse := CreatedResponse
	expectedSecondResponse := NoContentResponse
	responses := Responses{expectedFirstResponse, expectedSecondResponse}

	match := newMultipleResponsesMatch(request, responses)

	firstResponse, err := match.NextResponse(&http.Request{})
	s.NoError(err)
	s.Equal(expectedFirstResponse, firstResponse)

	secondResponse, err := match.NextResponse(&http.Request{})
	s.NoError(err)
	s.Equal(expectedSecondResponse, secondResponse)

	_, err = match.NextResponse(&http.Request{})
	s.ErrorIs(err, errNoNextResponseFound)
}
