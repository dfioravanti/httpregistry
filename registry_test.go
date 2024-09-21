package httpregistry_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/dfioravanti/httpregistry"
)

func (s *TestSuite) TestMockMatchesOnRouteWorks() {
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
			registry := httpregistry.NewRegistry()
			registry.AddSimpleRequest(tc.matchPath, http.MethodGet)

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

func (s *TestSuite) TestMockMatchesOnMethodWorks() {

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

			registry := httpregistry.NewRegistry()
			registry.AddSimpleRequest(path, tc.method)

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

func (s *TestSuite) TestMockMatchesOnHeader() {

	testCases := []struct {
		headers map[string]string
	}{
		{map[string]string{"Accept": "foo"}},
		{map[string]string{"Authorization": "Basic bar", "Accept": "foo"}},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("match %+v", tc.headers)
		path := "/test"
		client := http.Client{}

		s.Run(name, func() {

			registry := httpregistry.NewRegistry()
			registry.AddRequest(
				httpregistry.NewRequest(path, http.MethodGet, httpregistry.WithRequestHeaders(tc.headers)),
			)

			server := registry.GetServer(s.T())
			defer server.Close()

			request, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
			for k, v := range tc.headers {
				request.Header.Add(k, v)
			}
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

	registry := httpregistry.NewRegistry()
	mockRequest := httpregistry.NewRequest(path, http.MethodPost)
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

func (s *TestSuite) TestAccessTheRequestBodyByUrlWorks() {

	path := "/users"

	expectedBody := []byte(`
	{
		user_id: 10,
		foo: "bar",
		logged_in: true
	}
	`)

	registry := httpregistry.NewRegistry()
	registry.AddSimpleRequest(path, http.MethodPost)

	server := registry.GetServer(s.T())
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesUrl(path)
	s.Equal(1, len(matchingRequests))

	bodyBytes, err := io.ReadAll(matchingRequests[0].Body)
	s.NoError(err)

	s.Equal(expectedBody, bodyBytes)
}

func (s *TestSuite) TestAccessTheRequestBodyByUrlAndMethodWorks() {

	path := "/users"

	expectedBody := []byte(`
	{
		user_id: 10,
		foo: "bar",
		logged_in: true
	}
	`)

	registry := httpregistry.NewRegistry()
	registry.AddSimpleRequest(path, http.MethodPost)

	server := registry.GetServer(s.T())
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesUrlAndMethod(path, http.MethodPost)
	s.Equal(1, len(matchingRequests))

	bodyBytes, err := io.ReadAll(matchingRequests[0].Body)
	s.NoError(err)

	s.Equal(expectedBody, bodyBytes)
}

func (s *TestSuite) TestAccessTheRequestBodyByUrlAndMethodDoesNotMatchCorrectly() {

	path := "/users"

	expectedBody := []byte(`
	{
		user_id: 10,
		foo: "bar",
		logged_in: true
	}
	`)

	registry := httpregistry.NewRegistry()
	registry.AddSimpleRequest(path, http.MethodPost)

	server := registry.GetServer(s.T())
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesUrlAndMethod(path, http.MethodGet)
	s.Equal(0, len(matchingRequests))
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

	registry := httpregistry.NewRegistry()
	registry.AddRequestWithResponse(
		httpregistry.NewRequest(path, http.MethodPost),
		httpregistry.NewResponse(http.StatusCreated, expectedBody),
	)

	server := registry.GetServer(s.T())
	defer server.Close()
	client := http.Client{}

	request, err := http.NewRequest(http.MethodPost, server.URL+path, nil)
	s.NoError(err)

	res, err := client.Do(request)
	s.NoError(err)

	s.Equal(http.StatusCreated, res.StatusCode)

	bodyBytes, err := io.ReadAll(res.Body)
	s.NoError(err)

	s.Equal(expectedBody, bodyBytes)
}

func (s *TestSuite) TestMultipleCallsToTheSameURLWorks() {
	path := "/users"

	expectedFirstUser := []byte(`
	{
		user_id: 10,
	}
	`)
	expectedSecondUser := []byte(`
	{
		user_id: 20,
	}
	`)

	registry := httpregistry.NewRegistry()
	registry.AddRequestWithResponses(
		httpregistry.NewRequest(path, http.MethodGet),
		httpregistry.NewResponse(http.StatusCreated, expectedFirstUser),
		httpregistry.NewResponse(http.StatusCreated, expectedSecondUser),
		httpregistry.NewResponse(http.StatusNotFound, nil),
	)

	server := registry.GetServer(s.T())
	defer server.Close()
	client := http.Client{}

	request, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
	s.NoError(err)
	res, err := client.Do(request)
	s.NoError(err)
	s.Equal(http.StatusCreated, res.StatusCode)
	bodyBytes, err := io.ReadAll(res.Body)
	s.NoError(err)
	s.Equal(expectedFirstUser, bodyBytes)

	request, err = http.NewRequest(http.MethodGet, server.URL+path, nil)
	s.NoError(err)
	res, err = client.Do(request)
	s.NoError(err)
	s.Equal(http.StatusCreated, res.StatusCode)
	bodyBytes, err = io.ReadAll(res.Body)
	s.NoError(err)
	s.Equal(expectedSecondUser, bodyBytes)

	request, err = http.NewRequest(http.MethodGet, server.URL+path, nil)
	s.NoError(err)
	res, err = client.Do(request)
	s.NoError(err)
	s.Equal(http.StatusNotFound, res.StatusCode)
}
