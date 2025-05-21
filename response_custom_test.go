package httpregistry

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

func (s *TestSuite) TestFunctionalResponse() {
	testCases := []struct {
		name         string
		request      Request
		response     CustomResponse
		expectedBody string
	}{
		{
			name: "extract user id works",
			request: NewRequest().
				WithURL("users/12/address"),
			response: NewCustomResponse(func(w http.ResponseWriter, r *http.Request) {
				regexUser := regexp.MustCompile(`/users/(?P<userID>.+)/address$`)
				if regexUser.MatchString(r.URL.Path) {
					matches := regexUser.FindStringSubmatch(r.URL.Path)
					userID := matches[regexUser.SubexpIndex("userID")]
					body := map[string]string{"user_id": userID}
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(&body)
					return
				}
			}),
			expectedBody: "{\"user_id\":\"12\"}\n",
		},
	}
	for _, tc := range testCases {
		client := http.Client{}

		s.Run(tc.name, func() {
			registry := NewRegistry(s.T())
			registry.AddRequestWithResponse(tc.request, tc.response)

			server := registry.GetServer()
			defer server.Close()

			request, err := http.NewRequest(tc.request.Method, server.URL+"/"+tc.request.URL, bytes.NewReader(tc.request.Body))
			s.NoError(err)

			res, err := client.Do(request)
			s.NoError(err)

			body, err := io.ReadAll(res.Body)
			s.NoError(err)

			s.Equal(res.StatusCode, 200)
			s.Equal(tc.expectedBody, string(body))
		})
	}
}
