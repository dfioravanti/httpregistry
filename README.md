# HTTPRegistry

HTTPRegistry is a tiny package designed to simplifying building configurable `httptest` servers. [httptest](https://pkg.go.dev/net/http/httptest) is an incredibly powerful package that can be used to test the behavior of code that makes http calls. Unfortunately when the chain of calls to be tested is complex, setting up a mock server gets complicated and full of boilerplate test. This library is defined to take care of the boilerplate and let you focus on what matters in your tests.

## Basic concepts

In a nutshell this library allows you to create a `registry` on which `responses` to `requests` can be registered. Then this registry can be used to instantiate a httptest server that can respond to http requests. The library then takes care of checking that all the responses are used and that not too many calls happen.


```go
package main

import (
	"net/http"
	"testing"

	"github.com/dfioravanti/httpregistry"
)

func TestHttpRegistryWorks(t *testing.T) {
    // 1. Create the registry and defer the check that all responses
    //    that we will create are used.
    //    t is used to fail the test if the deferred check fails.
	registry := httpregistry.NewRegistry(t)
	defer registry.CheckAllResponsesAreConsumed()

    // 2. Add request to the registry
	registry.AddMethodAndURL(http.MethodGet, "/users")

    // 3. Create the server
	server := registry.GetServer()
	defer server.Close()

    // 4. Make calls, since we registered a single call then if we would call "/users" again it would fail
	response, err := http.Get(server.URL + "/users")
	if err != nil {
		t.Errorf("executing request failed: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %v", response.StatusCode)
	}
}
```
### Requests/Responses

The library provides various helper functions to make the process of attaching a response to a request easier. In the most general form it uses two types
* `Request` that defines the expected response to match
* `Response` that defines the response to be returned if a match happens

a request can be matched via the following

* A (method, exact path) combination, like `GET /users`
* A (method, regex) combination, like `GET /users/*`

plus additionally other constrains can be places on the matching

* It can be requested the request contains some headers like `Accept: text/html`.

Once a request is matched the corresponding `Response` is used to determine what the server should return. Currently the library allows to set

* Status code
* Body
* Headers

### Custom responses

Sometimes the standard `Response` from the package is not enough, suppose that you want to return a different value depending on the request, so for example you want to match an ID in the path or something similar. This is not possible with a `Response` since it does not allow to interact with the `http.Request` that is coming in. To solve this problem this package provides a `CustomResponse` type that allows you to interact with both the `http.Request` and the `http.ResponseWriter`.

```go
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/dfioravanti/httpregistry"
)

func TestCustomRequestWorks(t *testing.T) {
    // 1. Create the registry and defer the check that all responses
	//    that we will create are used.
	//    t is used to fail the test if the deferred check fails.
	registry := httpregistry.NewRegistry(t)
	defer registry.CheckAllResponsesAreConsumed()

	// 2. Create a CustomResponse, all functions that accept Response also accept CustomResponse
	mockResponse := httpregistry.NewCustomResponse(func(w http.ResponseWriter, r *http.Request) {
		regexUser := regexp.MustCompile(`/users/(?P<userID>.+)/address$`)
		if regexUser.MatchString(r.URL.Path) {
			matches := regexUser.FindStringSubmatch(r.URL.Path)
			userID := matches[regexUser.SubexpIndex("userID")]
			body := map[string]string{"user_id": userID}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(&body)
			return
		}
	})
	// Optional add a name to the CustomResponse so it is easier to debug, by default they get "custom response 1,2,3" as name
	mockResponse = mockResponse.WithName("match on user ID")
	registry.AddResponse(mockResponse)

	// 3. Create the server
	server := registry.GetServer()
	defer server.Close()

	// 4. Make calls and check assertions
	response, err := http.Get(server.URL + "/users/12/address")
	if err != nil {
		t.Errorf("executing request failed: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %v", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("reading body failed: %v", err)
	}

	expectedBody := "{\"user_id\":\"12\"}\n"
	if string(body) != expectedBody {
		t.Errorf("body does not match expected body")
	}
}
```

## How is a request selected

In case multiple requests match the incoming one then the first one, by order of registration, matching that still has unconsumed responses will be selected. So for example

```go
package main

import (
	"net/http"
	"testing"

	"github.com/dfioravanti/httpregistry"
)

func TestMultipleMatchingWorks(t *testing.T) {
	// 1. Create the registry and defer the check that all responses
	//    that we will create are used.
	//    t is used to fail the test if the deferred check fails.
	registry := httpregistry.NewRegistry(t)
	defer registry.CheckAllResponsesAreConsumed()

	// 2. Add requests to the registry
	registry.AddMethodAndURLWithStatusCode(
        http.MethodGet, "/users", http.StatusOK,
    )
	registry.AddMethodAndURLWithStatusCode(
        http.MethodGet, "/users", http.StatusNotFound,
    )

	// 3. Create the server
	server := registry.GetServer()
	defer server.Close()

	// 4. Make calls
	response, err := http.Get(server.URL + "/users")
	if err != nil {
		t.Errorf("executing request failed: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %v", response.StatusCode)
	}

	response, err = http.Get(server.URL + "/users")
	if err != nil {
		t.Errorf("executing request failed: %v", err)
	}
	if response.StatusCode != 404 {
		t.Errorf("unexpected status code %v", response.StatusCode)
	}
}
```
