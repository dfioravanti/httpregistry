package servermock_test

import (
	"net/http"

	"github.com/dfioravanti/servermock"
)

func (s *TestSuite) TestEqualityWorks() {
	r := servermock.NewRequest(
		servermock.WithRequestURL("/test"),
		servermock.WithRequestMethod(http.MethodPatch),
		servermock.WithRequestHeader("foo", "bar"),
	)

	s.True(r.Equal(r))
}
