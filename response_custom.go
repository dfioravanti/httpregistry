package httpregistry

import (
	"net/http"
)

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
	name string
	f    func(w http.ResponseWriter, r *http.Request)
}

// String marshal CustomResponse to string
func (res CustomResponse) String() string {
	return res.name
}

// serveResponse emits the response encoded in FunctionalResponse to w
func (res CustomResponse) serveResponse(w http.ResponseWriter, r *http.Request) {
	res.f(w, r)
}

// WithName allows to add a name to a FunctionalResponse so that it can be better identified when debugging.
// By the fault FunctionalResponse gets a sequential name that can be hard to identify if there are many of them
func (res CustomResponse) WithName(name string) CustomResponse {
	res.name = name
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
	return CustomResponse{
		name: "",
		f:    f,
	}
}
