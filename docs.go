/*
Package httpregistry provides multiple utilities that can be used to simplify the creation of /net/http/httptest
mock servers.
That package allows the creation of http servers that can be used to respond to actual http calls in tests.
This package aims at providing a nicer interface that should cover the most standard cases and attempts to hide away a layer of boilerplate.
For example it is normal to write test code like this

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/users" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))

with this package this can be simplified to

	registry := NewRegistry(t)
	registry.AddMethodAndURL("/users", http.MethodGet)
	ts := registry.GetServer()

Similarly this package tries to help with the harder task to test if a POST/PUT request actually happen to have the expected body/parameters.
With this library this can be done as

	registry := NewRegistry(t)
	registry.AddRequest(
		httpregistry.Request().
		WithURL("/users").
		WithMethod(http.MethodPost).
		WithJSONHeader().
		WithBody([]byte("{\"user\": \"John Schmidt\"}"))
	ts := registry.GetServer()

For more examples of what this package is capable of, refer to the README file.
*/
package httpregistry
