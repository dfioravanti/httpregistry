package httpregistry

import "fmt"

type whyMissed string

// These constants are used to represents why a match does not work and they are returned to the user as part of the error message so that errors can be debugged easily
const (
	pathDoesNotMatch   = whyMissed("the path does not match")
	methodDoesNotMatch = whyMissed("the method does not match")
	headerDoesNotMatch = whyMissed("the header does not match")
	outOfResponses     = whyMissed("the route matches but there was no response available")
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
	Request Request   `json:"request"`
	Why     whyMissed `json:"why"`
}

// newMiss creates a new miss and clones the match object to
// guarantee that it is not modified from outside changes
func newMiss(match match, why whyMissed) miss {
	return miss{match.Request(), why}
}

// String returns a human readable version of why the match could not happen
func (m miss) String() string {
	return fmt.Sprintf("%v missed because %v", m.Request, m.Why)
}
