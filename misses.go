package httpregistry

import "fmt"

type whyMissed string

// These constants are used to represents why a match does not work and they are returned to the user as part of the error message so that errors can be debugged easily
const (
	pathDoesNotMatch   = whyMissed("The path does not match")
	methodDoesNotMatch = whyMissed("The method does not match")
	headerDoesNotMatch = whyMissed("The header does not match")
)

// miss represents that the registry was not able to match a registered request with the current request that is coming in from the outside.
// This struct is used to communicate why a particular match cannot happen and it is designed to help the user to understand what went wrong.
//
// For example if the registered request is
//
//	POST /user
//
// but the incoming request is
//
//	PUT /user
//
// then the miss struct will communicate that the URL matches but the method does not
type miss struct {
	MissedMatch match
	Why         whyMissed
}

// String returns a human readable version of why the match could not happen
func (m miss) String() string {
	return fmt.Sprintf("%v missed %v", m.MissedMatch.Request(), m.Why)
}
