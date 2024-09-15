package servermock_test

import (
	"net/http"

	"github.com/dfioravanti/servermock"
)

func (s *TestSuite) TestFixedResponseMatchHasExpectedResponse() {
	request := servermock.NewRequest(
		servermock.WithRequestURL("/"),
		servermock.WithRequestMethod(http.MethodPost),
	)
	response := servermock.NoContentResponse

	match := servermock.NewFixedResponseMatch(request, response)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestFixedResponseRespondsForever() {
	request := servermock.NewRequest(
		servermock.WithRequestURL("/"),
		servermock.WithRequestMethod(http.MethodPost),
	)
	expectedResponse := servermock.NoContentResponse

	match := servermock.NewFixedResponseMatch(request, expectedResponse)
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
	request := servermock.NewRequest(
		servermock.WithRequestURL("/"),
		servermock.WithRequestMethod(http.MethodPost),
	)
	responses := servermock.Responses{servermock.NoContentResponse}

	match := servermock.NewMultipleResponsesMatch(request, responses)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestMultipleResponsesRespondsTheCorrectNumberOfTimes() {
	request := servermock.NewRequest(
		servermock.WithRequestURL("/"),
		servermock.WithRequestMethod(http.MethodPost),
	)

	expectedFirstResponse := servermock.CreatedResponse
	expectedSecondResponse := servermock.NoContentResponse
	responses := servermock.Responses{expectedFirstResponse, expectedSecondResponse}

	match := servermock.NewMultipleResponsesMatch(request, responses)

	firstResponse, err := match.NextResponse(&http.Request{})
	s.NoError(err)
	s.Equal(expectedFirstResponse, firstResponse)

	secondResponse, err := match.NextResponse(&http.Request{})
	s.NoError(err)
	s.Equal(expectedSecondResponse, secondResponse)

	_, err = match.NextResponse(&http.Request{})
	s.ErrorIs(err, servermock.ErrNoNextResponseFound)
}
