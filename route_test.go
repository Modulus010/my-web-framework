package wfw

import (
	"testing"
)

func TestParamsByName(t *testing.T) {
	params := Params{
		{Key: "id", Value: "123"},
		{Key: "name", Value: "test"},
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"id", "123"},
		{"name", "test"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := params.ByName(tt.key); got != tt.expected {
				t.Errorf("Params.ByName(%s) = %s; want %s", tt.key, got, tt.expected)
			}
		})
	}
}

func TestMethodRouteAddRoute(t *testing.T) {
	r := &methodRoute{}
	handlers := HandlersChain{func(c *Context) {}}

	r.addRoute("/users/:id", handlers)

	if len(r.nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(r.nodes))
	}

	node := r.nodes[0]
	expectedRegex := "^/users/([^/]+)$"
	if node.regex.String() != expectedRegex {
		t.Errorf("Expected regex %s, got %s", expectedRegex, node.regex.String())
	}

	if len(node.params) != 1 || node.params[0] != ":id" {
		t.Errorf("Expected params [\":id\"], got %v", node.params)
	}

	if len(node.handlers) != 1 || node.handlers[0] == nil {
		t.Errorf("Expected handlers, got %v", node.handlers)
	}
}

func TestMethodRouteGetValue(t *testing.T) {
	r := &methodRoute{}
	handlers := HandlersChain{func(c *Context) {}}

	r.addRoute("/users/:id", handlers)

	tests := []struct {
		path     string
		expected HandlersChain
		params   Params
	}{
		{"/users/123", handlers, Params{{Key: ":id", Value: "123"}}},
		{"/users/", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			gotHandlers, gotParams := r.getValue(tt.path)
			if len(gotHandlers) != len(tt.expected) {
				t.Errorf("Expected handlers %v, got %v", tt.expected, gotHandlers)
			}
			if len(gotParams) != len(tt.params) {
				t.Errorf("Expected params %v, got %v", tt.params, gotParams)
			}
		})
	}
}

func TestMethodRoutesGet(t *testing.T) {
	routes := methodRoutes{
		{method: "GET"},
		{method: "POST"},
	}

	tests := []struct {
		method   string
		expected *methodRoute
	}{
		{"GET", &routes[0]},
		{"POST", &routes[1]},
		{"PUT", nil},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			if got := routes.get(tt.method); got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}
