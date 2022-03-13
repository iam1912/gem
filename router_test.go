package gem

import (
	"testing"
)

func setupRouter() *Router {
	r := newRoter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/static/*filepath", nil)
	return r
}

type RouterTestCase struct {
	Method          string
	Pattern         string
	Key             string
	ExpectedPattern string
	ExpectedParams  string
	ExpectedNode    *node
}

func TestGetRouter(t *testing.T) {
	r := setupRouter()
	testCases := []RouterTestCase{
		{Method: "GET", Pattern: "/hello/world", Key: "name", ExpectedPattern: "/hello/:name", ExpectedParams: "world"},
		{Method: "GET", Pattern: "/static/index.html", Key: "filepath", ExpectedPattern: "/static/*filepath", ExpectedParams: "index.html"},
	}
	for i, testCase := range testCases {
		node, params := r.getRoute(testCase.Method, testCase.Pattern)
		if node != nil && node.pattern != testCase.ExpectedPattern {
			t.Errorf("TestGetRouter #%v: Expected get pattern %v but get %v\n", i+1, testCase.ExpectedPattern, node.pattern)
		}
		if params != nil && params[testCase.Key] != testCase.ExpectedParams {
			t.Errorf("TestGetRouter #%v: Expected get params %v but get %v\n", i+1, testCase.ExpectedParams, params[testCase.Key])
		}
	}
}
