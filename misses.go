package httpregistry

import "fmt"

type whyMissed string

const (
	pathDoesNotMatch   = whyMissed("The path does not match")
	methodDoesNotMatch = whyMissed("The method does not match")
	headersDoNotMatch  = whyMissed("The headers do not match")
)

type miss struct {
	MissedMatch Request
	Why         whyMissed
}

func (m miss) String() string {
	return fmt.Sprintf("%v missed %v", m.MissedMatch, m.Why)
}
