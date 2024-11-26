package httpregistry_test

import (
	"net/http"

	"github.com/dfioravanti/httpregistry"
)

func (s *TestSuite) TestEqualityWorks() {
	var testCases = []struct {
		name           string
		r1             httpregistry.Request
		r2             httpregistry.Request
		expectedResult bool
	}{
		{
			"equality works with default",
			httpregistry.NewRequest(),
			httpregistry.NewRequest(),
			true,
		},
		{
			"equality works with URL, method and header",
			httpregistry.NewRequest().WithURL("/test").WithMethod(http.MethodPatch).WithHeader("foo", "bar"),
			httpregistry.NewRequest().WithURL("/test").WithMethod(http.MethodPatch).WithHeader("foo", "bar"),
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.True(tc.r1.Equal(tc.r2))
		})
	}
}
