package wfw

import (
	"html/template"
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

func TestEngine_LoadHTMLGlob(t *testing.T) {
	engine := New()
	engine.LoadHTMLGlob("templates/*")

	if engine.HTMLTemplate == nil {
		t.Fatal("Expected HTMLTemplate to be loaded, but it was nil")
	}
}

func TestEngine_SetFuncMap(t *testing.T) {
	engine := New()
	funcMap := template.FuncMap{
		"testFunc": func() string { return "test" },
	}
	engine.SetFuncMap(funcMap)

	if engine.FuncMap == nil {
		t.Fatal("Expected FuncMap to be set, but it was nil")
	}

	if _, exists := engine.FuncMap["testFunc"]; !exists {
		t.Fatal("Expected FuncMap to contain 'testFunc', but it did not")
	}
}
