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
			registry := httpregistry.NewRegistry(s.T())
			registry.AddSimpleRequest(http.MethodGet, tc.matchPath)

			server := registry.GetServer()
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
			registry := httpregistry.NewRegistry(s.T())
			registry.AddSimpleRequest(tc.method, path)

			server := registry.GetServer()
			defer server.Close()

			request, err := http.NewRequest(tc.method, server.URL+path, nil)
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			s.Equal(http.StatusOK, res.StatusCode)
		})
	}
}

func (s *TestSuite) TestAddSimpleRequestWithStatusCodeWorks() {
	testCases := []struct {
		statusCode int
	}{
		{http.StatusOK},
		{http.StatusAccepted},
		{http.StatusNoContent},
		{http.StatusBadRequest},
		{http.StatusUnauthorized},
		{http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("returning %v", tc.statusCode)
		method := http.MethodGet
		path := "/test"
		client := http.Client{}

		s.Run(name, func() {
			registry := httpregistry.NewRegistry(s.T())
			registry.AddSimpleRequestWithStatusCode(method, path, tc.statusCode)

			server := registry.GetServer()
			defer server.Close()

			request, err := http.NewRequest(method, server.URL+path, nil)
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			s.Equal(tc.statusCode, res.StatusCode)
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
			registry := httpregistry.NewRegistry(s.T())
			registry.AddRequest(httpregistry.NewRequest(http.MethodGet, path, httpregistry.WithRequestHeaders(tc.headers)))

			server := registry.GetServer()
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

	registry := httpregistry.NewRegistry(s.T())
	mockRequest := httpregistry.NewRequest(http.MethodPost, path)
	registry.AddRequest(mockRequest)

	server := registry.GetServer()
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesForRequest(mockRequest)
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

	registry := httpregistry.NewRegistry(s.T())
	registry.AddSimpleRequest(http.MethodPost, path)

	server := registry.GetServer()
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesForURL(path)
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

	registry := httpregistry.NewRegistry(s.T())
	registry.AddSimpleRequest(http.MethodPost, path)

	server := registry.GetServer()
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesURLAndMethod(path, http.MethodPost)
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

	registry := httpregistry.NewRegistry(s.T())
	registry.AddSimpleRequest(http.MethodPost, path)

	server := registry.GetServer()
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewBuffer(expectedBody))
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	matchingRequests := registry.GetMatchesURLAndMethod(path, http.MethodGet)
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

	registry := httpregistry.NewRegistry(s.T())
	registry.AddRequestWithResponse(
		httpregistry.NewRequest(http.MethodPost, path),
		httpregistry.NewResponse(http.StatusCreated, expectedBody),
	)

	server := registry.GetServer()
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

	registry := httpregistry.NewRegistry(s.T())
	registry.AddRequestWithResponses(
		httpregistry.NewRequest(http.MethodGet, path),
		httpregistry.NewResponse(http.StatusCreated, expectedFirstUser),
		httpregistry.NewResponse(http.StatusCreated, expectedSecondUser),
		httpregistry.NewResponse(http.StatusNotFound, nil),
	)

	server := registry.GetServer()
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
