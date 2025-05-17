package httpregistry

import (
	"fmt"
	"io"
	"net/http"
)

type typeMarshable struct {
	A int `json:"a"`
}

func (t typeMarshable) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"a\": %d}", t.A)), nil
}

func (s *TestSuite) TestCustomResponseReturnsExpectedHeadersAndBody() {
	testCases := []struct {
		name            string
		response        Response
		expectedBody    string
		expectedHeaders http.Header
	}{
		{
			name:            "json response works",
			response:        NewResponse().WithJSONBody("{\"a\": 10}"),
			expectedBody:    "{\"a\": 10}",
			expectedHeaders: map[string][]string{"Content-Type": {"application/json"}},
		},
		{
			name:            "json marshaler response works",
			response:        NewResponse().WithJSONMarshalerBody(typeMarshable{A: 10}),
			expectedBody:    "{\"a\": 10}",
			expectedHeaders: map[string][]string{"Content-Type": {"application/json"}},
		},
	}
	for _, tc := range testCases {
		method := http.MethodGet
		path := "/test"
		request := NewRequest().WithMethod(method).WithURL(path)
		client := http.Client{}

		s.Run(tc.name, func() {
			registry := NewRegistry(s.T())
			registry.AddRequestWithResponse(request, tc.response)

			server := registry.GetServer()
			defer server.Close()

			request, err := http.NewRequest(method, server.URL+path, nil)
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			body, err := io.ReadAll(res.Body)
			s.NoError(err)

			s.Equal(res.StatusCode, 200)
			s.Equal(tc.expectedBody, string(body))
			s.Subset(res.Header, tc.expectedHeaders)
		})
	}
}
