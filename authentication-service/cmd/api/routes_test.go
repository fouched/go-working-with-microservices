package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"testing"
)

func Test_routes_exist(t *testing.T) {
	testApp := Config{}

	// returns a handler
	testRoutes := testApp.routes()
	// cast it to router
	chiRoutes := testRoutes.(chi.Router)

	// add all the routes that needs to be tested for existence - we only have one
	routes := []string{"/authenticate"}

	for _, route := range routes {
		routeExits(t, chiRoutes, route)
	}
}

func routeExits(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(handler2 http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}
