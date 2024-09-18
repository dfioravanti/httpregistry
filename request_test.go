package servermock_test

import (
	"net/http"

	"github.com/dfioravanti/servermock"
)

func (s *TestSuite) TestEqualityWorks() {
	r := servermock.NewRequest(
		"test",
		http.MethodPatch,
		servermock.WithRequestHeader("foo", "bar"),
	)

	s.True(r.Equal(r))
}
