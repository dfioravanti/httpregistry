package httpregistry_test

import (
	"net/http"

	"github.com/dfioravanti/httpregistry"
)

func (s *TestSuite) TestEqualityWorks() {
	r := httpregistry.NewRequest(
		"test",
		http.MethodPatch,
		httpregistry.WithRequestHeader("foo", "bar"),
	)

	s.True(r.Equal(r))
}
