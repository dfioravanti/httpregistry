package servermock_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/dfioravanti/servermock"
)

func (s *TestSuite) TestMockMatchesOnRoute() {
	testCases := []struct {
		matchPath  string
		calledPath string
	}{
		{"/test", "/test"},
		{"/users/(.*?)", "/users/123"},
		{"/foo\\?bar=(.*?)", "/foo?bar=10"},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("match %s with %s", tc.matchPath, tc.calledPath)
		client := http.Client{}

		s.Run(name, func() {
			registry := servermock.NewRegistry()
			registry.AddSimpleRequest(http.MethodGet, tc.matchPath)

			server := registry.GetServer(s.T())
			defer server.Close()

			request, err := http.NewRequest(http.MethodGet, server.URL+tc.calledPath, nil)
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			s.Equal(http.StatusOK, res.StatusCode)
		})
	}
}

func (s *TestSuite) TestMockMatchesOnMethod() {

	testCases := []struct {
		method string
	}{
		{http.MethodGet},
		{http.MethodHead},
		{http.MethodPost},
		{http.MethodPut},
		{http.MethodPatch},
		{http.MethodDelete},
		{http.MethodConnect},
		{http.MethodOptions},
		{http.MethodTrace},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("match %s", tc.method)
		path := "/test"
		client := http.Client{}

		s.Run(name, func() {

			registry := servermock.NewRegistry()
			registry.AddSimpleRequest(tc.method, path)

			server := registry.GetServer(s.T())
			defer server.Close()

			request, err := http.NewRequest(tc.method, server.URL+path, nil)
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			s.Equal(http.StatusOK, res.StatusCode)
		})
	}
}

func (s *TestSuite) TestAccessTheRequestBodyWorks() {

	path := "/users"

	expectedBody := []byte(`
	{
		user_id: 10,
		foo: "bar",
		logged_in: true
	}
	`)

	registry := servermock.NewRegistry()
	mockRequest := servermock.NewRequest(
		servermock.WithRequestMethod(http.MethodPost),
		servermock.WithRequestURL(path),
	)
	registry.AddRequest(mockRequest)

	server := registry.GetServer(s.T())
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesPerRequest(mockRequest)
	s.Equal(1, len(matchingRequests))

	bodyBytes, err := io.ReadAll(matchingRequests[0].Body)
	s.NoError(err)

	s.Equal(expectedBody, bodyBytes)
}

func (s *TestSuite) TestMockResponseWithBody() {

	path := "/users"

	expectedBody := []byte(`
	{
		user_id: 10,
		foo: "bar",
		logged_in: true
	}
	`)

	mock := servermock.NewRegistry()
	mock.AddRequestWithResponse(
		servermock.NewRequest(
			servermock.WithRequestMethod(http.MethodPost),
			servermock.WithRequestURL(path),
		),
		servermock.NewResponse(
			servermock.WithResponseStatus(http.StatusCreated),
			servermock.WithResponseJSONBody(expectedBody),
		),
	)

	server := mock.GetServer(s.T())
	defer server.Close()

	request, err := http.NewRequest(http.MethodPost, server.URL+path, nil)
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(request)
	s.NoError(err)

	s.Equal(http.StatusCreated, res.StatusCode)

	bodyBytes, err := io.ReadAll(res.Body)
	s.NoError(err)

	s.Equal(expectedBody, bodyBytes)
}
