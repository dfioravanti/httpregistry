package httpregistry

import (
	"net/http"
)

func (s *TestSuite) TestMultipleResponsesHasExpectedResponse() {
	request := NewRequest().WithMethod(http.MethodPost).WithURL("/")
	responses := mockResponses{NoContentResponse}

	match := newConsumableResponsesMatch(request, responses)

	s.Equal(request, match.Request())
}

func (s *TestSuite) TestMultipleResponsesRespondsTheCorrectNumberOfTimes() {
	request := NewRequest().WithMethod(http.MethodPost).WithURL("/")

	expectedFirstResponse := CreatedResponse
	expectedSecondResponse := NoContentResponse
	responses := mockResponses{expectedFirstResponse, expectedSecondResponse}

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
