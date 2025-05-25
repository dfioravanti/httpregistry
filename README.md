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

### Retrieving matching requests

Sometimes it is useful to retrieve the `http.Request` that matched a `HttpRegistry.Request`. For example to check if the body if a POST/PUT request is what we expect or similar. One can do that via three functions

- `registry.GetMatchesForRequest(request)`
- `registry.GetMatchesForURL(url)`
- `registry.GetMatchesURLAndMethod(url, method)`

### Infinite responses

A `Response` is consumed when a match happen, this is by design so that it is possible to test that the expected number of calls happens, but sometimes one does not really care about how many calls are made and just wants to mock a http call away. This is possible via `httpregistry.AddInfiniteResponse(response)`

```go
import (
	"net/http"
	"testing"

	"github.com/dfioravanti/httpregistry"
)

func TestAddInfiniteResponse(t *testing.T) {
	nbCalls := 100

	// 1. Create the registry and defer the check that all responses
	//    that we will create are used.
	//    t is used to fail the test if the deferred check fails.
	registry := httpregistry.NewRegistry(t)

	// 2. Add an AddInfiniteResponse,
	// 	  by default Response are consumed when they match but
	//    InfiniteResponse are not so they can be matched forever
	registry.AddInfiniteResponse(
		httpregistry.NewResponse(),
	)

	// 3. Create the server
	server := registry.GetServer()
	defer server.Close()
	client := http.Client{}

	// 4. Make calls and check assertions
	var urls []string
	for range nbCalls {
		url := generateRandomString(10)

		res, err := client.Get(server.URL + "/" + url)
		if err != nil {
			t.Errorf("executing request failed: %v", err)
		}

		if res.StatusCode != 200 {
			t.Errorf("unexpected status code %v", res.StatusCode)
		}

		urls = append(urls, url)
	}

	requests := registry.GetMatchesForRequest(httpregistry.DefaultRequest)
	if len(requests) != nbCalls {
		t.Errorf("the number of requests (%d) does not match the number of calls (%d)", len(requests), nbCalls)
	}

	for i, r := range requests {
		if r.URL.Path != "/"+urls[i] {
			t.Errorf("the request path (%s) does not match the expected path (%s)", r.URL.Path, "/"+urls[i])
		}
	}
}
```

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

## Investigate failed tests

The library tries to help as much as possible in debugging why a test has failed. To achieve this it will
1. Fail a test if
   1. It is impossible to reply to a request
   2. `registry.CheckAllResponsesAreConsumed()` is called but not all the requests are consumed
2. In case if it is impossible to reply to a request it will report in the body of the response why it failed
3. Provide a `httpregistry.NewMockTestingT()` that can be passed in place of `*testing.T` so that test failures can be better analyzed

### Investigate if a test fails to consume all requests
```go
package main

import (
	"net/http"
	"testing"

	"github.com/dfioravanti/httpregistry"
)

func TestHowToInvestigateFailingTest(t *testing.T) {
	// 1. Setup a mock for testing.T that we control and can access later
	mockT := httpregistry.NewMockTestingT()

	// 2. Setup registry and the requests
	registry := httpregistry.NewRegistry(mockT)
	registry.AddMethodAndURL(http.MethodGet, "/foo")
	registry.AddMethodAndURL(http.MethodDelete, "/bar")
	registry.AddRequestWithResponse(
		httpregistry.DefaultRequest,
		httpregistry.NewResponse().WithName("My beautiful response"),
	)
	registry.AddRequestWithResponse(
		httpregistry.DefaultRequest,
		httpregistry.NewCustomResponse(func(w http.ResponseWriter, r *http.Request) {}).WithName("My beautiful custom response"),
	)

	// 3. No call happens but we assert that all calls were consumed
	registry.CheckAllResponsesAreConsumed()

	// 4. Let us check that mockT contains useful information
	if len(mockT.Messages) != 4 {
		t.Errorf("There should be 4 uncalled request but I found only %d", len(mockT.Messages))
	}

	if !slices.Contains(mockT.Messages, "request mock request #1 has httpregistry.OkResponse as unused response") {
		t.Error("request mock request #1 has httpregistry.OkResponse as unused response should be in the slice but it is not")
	}
	if !slices.Contains(mockT.Messages, "request mock request #2 has httpregistry.OkResponse as unused response") {
		t.Error("request mock request #2 has httpregistry.OkResponse as unused response should be in the slice but it is not")
	}
	if !slices.Contains(mockT.Messages, "request httpregistry.DefaultRequest has My beautiful response as unused response") {
		t.Error("request httpregistry.DefaultRequest has My beautiful response as unused response should be in the slice but it is not")
	}
	if !slices.Contains(mockT.Messages, "request httpregistry.DefaultRequest has My beautiful custom response as unused response") {
		t.Error("request httpregistry.DefaultRequest has My beautiful custom response as unused response should be in the slice but it is not")
	}
}
```
### Investigate why a test fails when calling

```go
package main

import (
	"net/http"
	"testing"

	"github.com/dfioravanti/httpregistry"
)

func TestWeCanInvestigateWhyATestFails(t *testing.T) {
	// 1. Setup a mock for testing.T that we control and can access later
	mockT := httpregistry.NewMockTestingT()

	// 2. Setup registry and the requests
	registry := httpregistry.NewRegistry(mockT)
	registry.AddMethodAndURL(http.MethodGet, "/foo")

	// 3. Call Twice a route with only one response
	url := registry.GetServer().URL
	client := http.Client{}

	// 3a. First call works
	firstResponse, err := client.Get(url + "/foo")
	if err != nil {
		t.Errorf("Unexpected error in first request: %s", err)
	}
	if firstResponse.StatusCode != 200 {
		t.Errorf("Unexpected status code for first response, I was expecting 200 we got %d", firstResponse.StatusCode)
	}

	// 3b. Second call fails
	secondResponse, err := client.Get(url + "/foo")
	if err != nil {
		t.Errorf("Unexpected error in second request: %s", err)
	}
	if secondResponse.StatusCode != 500 {
		t.Errorf("Unexpected status code for second response, I was expecting 500 we got %d", secondResponse.StatusCode)
	}

	// 4. The test was failed
	if mockT.HasFailed != true {
		t.Errorf("mockT.HasFailed should be true, but it was %t", mockT.HasFailed)
	}

	// 5. The body of the call tells us why it failed
	bodyBytes, err := io.ReadAll(secondResponse.Body)
	if err != nil {
		t.Errorf("Decoding second response body failed: %s", err)
	}
	body := string(bodyBytes)
	if body != "mock request #1 missed because the route matches but there was no response available" {
		t.Errorf("was expecting \"mock request #1 missed because the route matches but there was no response available\", got: %s", body)
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
