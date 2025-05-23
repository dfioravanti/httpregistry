package httpregistry_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/dfioravanti/httpregistry"
)

func (s *TestSuite) TestMockMatchesOnURLWorks() {
	testCases := []struct {
		matchPath    string
		calledPath   string
		calledMethod string
	}{
		{"/test", "/test", http.MethodGet},
		{"/users/(.*?)", "/users/123", http.MethodDelete},
		{"/foo\\?bar=(.*?)", "/foo?bar=10", http.MethodPut},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("match %s with %s", tc.matchPath, tc.calledPath)
		client := http.Client{}

		s.Run(name, func() {
			registry := httpregistry.NewRegistry(s.T())
			registry.AddURL(tc.matchPath)

			server := registry.GetServer()
			defer server.Close()

			request, err := http.NewRequest(tc.calledMethod, server.URL+tc.calledPath, nil)
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			s.Equal(http.StatusOK, res.StatusCode)
		})
	}
}

func (s *TestSuite) TestMockMatchesOnRouteAndMethodWorks() {
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
			registry.AddMethodAndURL(http.MethodGet, tc.matchPath)

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

func (s *TestSuite) TestDefaultsMatchesEverything() {
	testCases := []struct {
		url    string
		method string
	}{
		{"/foo", http.MethodGet},
		{"/bar", http.MethodHead},
		{"test", http.MethodPost},
		{"/foo", http.MethodPut},
		{"/bar", http.MethodPatch},
		{"test", http.MethodDelete},
		{"/foo", http.MethodConnect},
		{"/bar", http.MethodOptions},
		{"test", http.MethodTrace},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("match %s", tc.method)
		path := "/test"
		client := http.Client{}

		s.Run(name, func() {
			registry := httpregistry.NewRegistry(s.T())
			registry.AddRequest(httpregistry.NewRequest())

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
			registry.AddMethod(tc.method)

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

func (s *TestSuite) TestAddBodyWorks() {
	testCases := []struct {
		body []byte
	}{
		{[]byte("body")},
		{[]byte("{\"message\": \"Hello, firstName()! Your order number is: #int(1,100)\"}")},
		{[]byte("<root></root>")},
		{[]byte{187, 163, 35, 30}},
	}
	for i, tc := range testCases {
		name := fmt.Sprintf("case %v", i)
		path := "/test"
		client := http.Client{}

		s.Run(name, func() {
			registry := httpregistry.NewRegistry(s.T())
			registry.AddBody(tc.body)

			server := registry.GetServer()
			defer server.Close()

			request, err := http.NewRequest(http.MethodGet, server.URL+path, bytes.NewReader(tc.body))
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			s.Equal(http.StatusOK, res.StatusCode)
		})
	}
}

func (s *TestSuite) TestAddMethodAndURLWithStatusCodeWorks() {
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
			registry.AddMethodAndURLWithStatusCode(method, path, tc.statusCode)

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
			registry.AddRequest(httpregistry.NewRequest().WithMethod(http.MethodGet).WithURL(path).WithHeaders(tc.headers))

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
	mockRequest := httpregistry.NewRequest().WithMethod(http.MethodPost).WithURL(path)
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

	// we check that the body is not consumed but it can be accessed again
	matchingRequests = registry.GetMatchesForRequest(mockRequest)
	s.Equal(1, len(matchingRequests))

	bodyBytes, err = io.ReadAll(matchingRequests[0].Body)
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
	registry.AddMethodAndURL(http.MethodPost, path)

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

	// we check that the body is not consumed but it can be accessed again
	matchingRequests = registry.GetMatchesForURL(path)
	s.Equal(1, len(matchingRequests))

	bodyBytes, err = io.ReadAll(matchingRequests[0].Body)
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
	registry.AddMethodAndURL(http.MethodPost, path)

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

	// we check that the body is not consumed but it can be accessed again
	matchingRequests = registry.GetMatchesURLAndMethod(path, http.MethodPost)
	s.Equal(1, len(matchingRequests))

	bodyBytes, err = io.ReadAll(matchingRequests[0].Body)
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
	registry.AddMethodAndURL(http.MethodPost, path)

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
		httpregistry.NewRequest().WithMethod(http.MethodPost).WithURL(path),
		httpregistry.NewResponse().WithStatus(http.StatusCreated).WithBody(expectedBody),
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
		httpregistry.NewRequest().WithMethod(http.MethodGet).WithURL(path),
		httpregistry.NewResponse().WithStatus(http.StatusCreated).WithBody(expectedFirstUser),
		httpregistry.NewResponse().WithStatus(http.StatusCreated).WithBody(expectedSecondUser),
		httpregistry.NewResponse().WithStatus(http.StatusNotFound),
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

	registry.CheckAllResponsesAreConsumed()
}

func (s *TestSuite) TestTestUsingAllRoutesWorks() {
	registry := httpregistry.NewRegistry(s.T())
	registry.AddMethodAndURL(http.MethodGet, "/foo")
	registry.AddMethodAndURL(http.MethodGet, "/bar")

	url := registry.GetServer().URL
	client := http.Client{}
	_, _ = client.Get(url + "/foo")
	_, _ = client.Get(url + "/bar")
	registry.CheckAllResponsesAreConsumed()
}

func (s *TestSuite) TestUncalledRoutesTriggerAFailure() {
	mockT := &httpregistry.MockTestingT{}

	registry := httpregistry.NewRegistry(mockT)
	registry.AddMethodAndURL(http.MethodGet, "/foo")
	registry.AddMethodAndURL(http.MethodDelete, "/bar")

	registry.CheckAllResponsesAreConsumed()

	s.Equal(len(mockT.Messages), 2)
	s.Contains(mockT.Messages, "request {\"url\":\"/foo\",\"method\":\"GET\"} has {\"status_code\":200} as unused response")
	s.Contains(mockT.Messages, "request {\"url\":\"/bar\",\"method\":\"DELETE\"} has {\"status_code\":200} as unused response")
}

func (s *TestSuite) TestCallingTooMayTimesFails() {
	mockT := &httpregistry.MockTestingT{}

	registry := httpregistry.NewRegistry(mockT)
	registry.AddMethodAndURL(http.MethodGet, "/foo")

	url := registry.GetServer().URL
	client := http.Client{}
	_, _ = client.Get(url + "/foo")
	response, err := client.Get(url + "/foo")
	s.NoError(err)

	bodyBytes, err := io.ReadAll(response.Body)
	s.NoError(err)
	body := string(bodyBytes)

	s.True(mockT.HasFailed)
	s.Equal(body, "[{\"request\":{\"url\":\"/foo\",\"method\":\"GET\"},\"why\":\"The route matches but there was no response available\"}]")
}

func (s *TestSuite) TestMatchArbitraryRequests() {
	testCases := []struct {
		name          string
		request       httpregistry.Request
		methodToCall  string
		pathToCall    string
		bodyToCall    []byte
		headersToCall http.Header
	}{
		{
			name: "match url",
			request: httpregistry.NewRequest().
				WithURL("/foo"),
			methodToCall: http.MethodGet,
			pathToCall:   "/foo",
		},
		{
			name: "match json body",
			request: httpregistry.NewRequest().
				WithURL("/foo").
				WithMethod(http.MethodPost).
				WithJSONBody(map[string]int{"foo": 10, "bar": 20}),
			methodToCall: http.MethodPost,
			pathToCall:   "/foo",
			bodyToCall:   mustMarshalJSON(map[string]int{"foo": 10, "bar": 20}),
		},
	}
	for _, tc := range testCases {
		client := http.Client{}

		s.Run(tc.name, func() {
			registry := httpregistry.NewRegistry(s.T())
			registry.AddRequest(tc.request)

			server := registry.GetServer()
			defer server.Close()

			b := bytes.NewReader(tc.bodyToCall)
			request, err := http.NewRequest(tc.methodToCall, server.URL+tc.pathToCall, b)
			s.NoError(err)

			for k, hs := range tc.headersToCall {
				for _, h := range hs {
					request.Header.Add(k, h)
				}
			}

			res, err := client.Do(request)
			s.NoError(err)

			s.Equal(http.StatusOK, res.StatusCode)
		})
	}
}
