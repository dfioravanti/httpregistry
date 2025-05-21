package httpregistry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// this is a bit of hacky thing to have a default name
var counter = 0

// CustomResponse allows the user to define a custom made response to any request.
// In particular it allows to define responses that are functions of the request
//
// for example
//
//	func(w http.ResponseWriter, r *http.Request) {
//		regexUser := regexp.MustCompile(`/users/(?P<userID>.+)/address$`)
//		if regexUser.MatchString(r.URL.Path) {
//			matches := regexUser.FindStringSubmatch(r.URL.Path)
//			userID := matches[regexUser.SubexpIndex("userID")]
//			body := map[string]string{"user_id": userID}
//			w.Header().Set("Content-Type", "application/json")
//			_ = json.NewEncoder(w).Encode(&body)
//			return
//		}
//	}
type CustomResponse struct {
	Name string `json:"name,omitempty"`
	f    func(w http.ResponseWriter, r *http.Request)
}

// String marshal FunctionalResponse to string
func (res CustomResponse) String() string {
	bytes, err := json.Marshal(res)
	if err != nil {
		panic("cannot marshal request")
	}
	return string(bytes)
}

// createResponse emits the response encoded in FunctionalResponse to w
func (res CustomResponse) createResponse(w http.ResponseWriter, r *http.Request) {
	res.f(w, r)
}

// WithName allows to add a name to a FunctionalResponse so that it can be better identified when debugging.
// By the fault FunctionalResponse gets a sequential name that can be hard to identify if there are many of them
func (res CustomResponse) WithName(name string) CustomResponse {
	res.Name = name
	return res
}

// NewCustomResponse creates a new FunctionalResponse.
// A FunctionalResponse allows the user to define a custom made response to any request.
// In particular it allows to define responses that are functions of the request
//
// for example
//
//	func(w http.ResponseWriter, r *http.Request) {
//		regexUser := regexp.MustCompile(`/users/(?P<userID>.+)/address$`)
//		if regexUser.MatchString(r.URL.Path) {
//			matches := regexUser.FindStringSubmatch(r.URL.Path)
//			userID := matches[regexUser.SubexpIndex("userID")]
//			body := map[string]string{"user_id": userID}
//			w.Header().Set("Content-Type", "application/json")
//			_ = json.NewEncoder(w).Encode(&body)
//			return
//		}
//	}
func NewCustomResponse(f func(w http.ResponseWriter, r *http.Request)) CustomResponse {
	res := CustomResponse{Name: fmt.Sprintf("Custom response %d", counter), f: f}
	counter++
	return res
}
