package httpregistry_test

import (
	"net/http"

	"github.com/dfioravanti/httpregistry"
)

func (s *TestSuite) TestFixedResponseMatchHasExpectedResponse() {
	request := httpregistry.NewRequest("/", http.MethodPost)
	response := httpregistry.NoContentResponse

	match := httpregistry.NewFixedResponseMatch(request, response)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestFixedResponseRespondsForever() {
	request := httpregistry.NewRequest("/", http.MethodPost)
	expectedResponse := httpregistry.NoContentResponse

	match := httpregistry.NewFixedResponseMatch(request, expectedResponse)
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
	request := httpregistry.NewRequest("/", http.MethodPost)
	responses := httpregistry.Responses{httpregistry.NoContentResponse}

	match := httpregistry.NewMultipleResponsesMatch(request, responses)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestMultipleResponsesRespondsTheCorrectNumberOfTimes() {
	request := httpregistry.NewRequest("/", http.MethodPost)

	expectedFirstResponse := httpregistry.CreatedResponse
	expectedSecondResponse := httpregistry.NoContentResponse
	responses := httpregistry.Responses{expectedFirstResponse, expectedSecondResponse}

	match := httpregistry.NewMultipleResponsesMatch(request, responses)

	firstResponse, err := match.NextResponse(&http.Request{})
	s.NoError(err)
	s.Equal(expectedFirstResponse, firstResponse)

	secondResponse, err := match.NextResponse(&http.Request{})
	s.NoError(err)
	s.Equal(expectedSecondResponse, secondResponse)

	_, err = match.NextResponse(&http.Request{})
	s.ErrorIs(err, httpregistry.ErrNoNextResponseFound)
}
