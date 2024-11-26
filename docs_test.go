package httpregistry_test

import (
	"net/http"
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
