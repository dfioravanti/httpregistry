package httpregistry

import "fmt"

type whyMissed string

const (
	pathDoesNotMatch   = whyMissed("The path does not match")
	methodDoesNotMatch = whyMissed("The method does not match")
)

type Miss struct {
	MissedMatch Request
	Why         whyMissed
}

func (m Miss) String() string {
	return fmt.Sprintf("%v missed %v", m.MissedMatch, m.Why)
}
