package wfw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_Run(t *testing.T) {
	engine := New()
	engine.addRoute("GET", "/", HandlersChain{func(c *Context) {
		c.Writer.WriteHeader(http.StatusOK)
	}})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}

func TestEngine_NotFound(t *testing.T) {
	engine := New()

	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status code %d, but got %d", http.StatusNotFound, w.Code)
	}
}
