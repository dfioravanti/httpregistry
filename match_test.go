package httpregistry

import (
	"net/http"
)

func (s *TestSuite) TestMultipleResponsesHasExpectedResponse() {
	request := NewRequest(http.MethodPost, "/")
	responses := Responses{NoContentResponse}

	match := newConsumableResponsesMatch(request, responses)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestMultipleResponsesRespondsTheCorrectNumberOfTimes() {
	request := NewRequest(http.MethodPost, "/")

	expectedFirstResponse := CreatedResponse
	expectedSecondResponse := NoContentResponse
	responses := Responses{expectedFirstResponse, expectedSecondResponse}

	match := newConsumableResponsesMatch(request, responses)

	firstResponse, err := match.NextResponse()
	s.NoError(err)
	s.Equal(expectedFirstResponse, firstResponse)

	secondResponse, err := match.NextResponse()
	s.NoError(err)
	s.Equal(expectedSecondResponse, secondResponse)

	_, err = match.NextResponse()
	s.ErrorIs(err, errNoNextResponseFound)
}
