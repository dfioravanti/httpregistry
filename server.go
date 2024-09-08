package servermock

import (
	"net/http"
	"net/http/httptest"
)

type Match struct {
	request       Request
	response      Response
	NumberOfCalls int
}

type whyMissed string

const (
	pathDoesNotMatch   = whyMissed("The path does not match")
	methodDoesNotMatch = whyMissed("The method does not match")
)

type Miss struct {
	MissedMatch Request
	Why         whyMissed
}

type Mock struct {
	matches []Match
	misses  []Miss
}

func NewMock() *Mock {
	mock := Mock{}
	return &mock
}

func (m *Mock) AddRequest(request Request) {
	m.matches = append(m.matches, Match{
		request:  request,
		response: Response200,
	})
}

func (m *Mock) GetServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathToMatch := r.URL.String()
		methodToMatch := r.Method

		for _, possibleMatch := range m.matches {
			if possibleMatch.request.urlAsRegex.MatchString(pathToMatch) {
				if methodToMatch == r.Method {
					for k, v := range possibleMatch.response.headers {
						w.Header().Add(k, v)
					}
					w.WriteHeader(possibleMatch.response.status)
					w.Write(possibleMatch.response.body)
					return
				} else {
					miss := Miss{
						MissedMatch: possibleMatch.request,
						Why:         methodDoesNotMatch,
					}
					m.misses = append(m.misses, miss)
				}
			} else {
				miss := Miss{
					MissedMatch: possibleMatch.request,
					Why:         pathDoesNotMatch,
				}
				m.misses = append(m.misses, miss)
			}
		}

		// If nothing matches we return 500 because we assume that something should match.
		// Moreover returning 500 allows to distinguish the case where returning 404 is intentional
		w.WriteHeader(http.StatusInternalServerError)
	}))
}
