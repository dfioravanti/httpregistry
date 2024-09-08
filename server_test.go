package servermock

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockMatchesOnRoute(t *testing.T) {
	t.Parallel()

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mock := NewMock()
			mock.AddRequest(NewRequest(
				WithRequestURL(tc.matchPath),
				WithRequestMethod(http.MethodGet),
			))

			server := mock.GetServer()
			defer server.Close()

			request, err := http.NewRequest(http.MethodGet, server.URL+tc.calledPath, nil)
			assert.NoError(t, err)

			res, err := client.Do(request)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}

func TestMockMatchesOnMethod(t *testing.T) {
	t.Parallel()

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mock := NewMock()
			mock.AddRequest(NewRequest(
				WithRequestURL(path),
				WithRequestMethod(tc.method),
			))

			server := mock.GetServer()
			defer server.Close()

			request, err := http.NewRequest(tc.method, server.URL+path, nil)
			assert.NoError(t, err)

			res, err := client.Do(request)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}
