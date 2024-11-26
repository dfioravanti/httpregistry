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

    // 4. Make calls
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
