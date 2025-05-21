package httpregistry_test

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/dfioravanti/httpregistry"
)

// This file contains all the code that is embedded in the readme.
// If anything here breaks remember to update the readme.

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

func TestMultipleMatchingWorks(t *testing.T) {
	// 1. Create the registry and defer the check that all responses
	//    that we will create are used.
	//    t is used to fail the test if the deferred check fails.
	registry := httpregistry.NewRegistry(t)
	defer registry.CheckAllResponsesAreConsumed()

	// 2. Add requests to the registry
	registry.AddMethodAndURLWithStatusCode(http.MethodGet, "/users", http.StatusOK)
	registry.AddMethodAndURLWithStatusCode(http.MethodGet, "/users", http.StatusNotFound)

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
