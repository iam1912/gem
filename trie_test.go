package gem

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TrieTestCase struct {
	Pattern         string
	ExpectedPattern string
	ExpectedNode    *node
	OperateType     string
}

func setupTrie() *node {
	initialization := []TrieTestCase{
		{Pattern: "/hello/name"},
		{Pattern: "/search/:name"},
		{Pattern: "/static/*filepath"},
	}
	node := &node{}
	for _, value := range initialization {
		parts := parsePattern(value.Pattern)
		node.insert(value.Pattern, parts, 0)
	}
	return node
}

func TestTrie(t *testing.T) {
	root := setupTrie()

	testCases := []TrieTestCase{
		{Pattern: "/hello/name", ExpectedPattern: "/hello/name", OperateType: "search"},
		{Pattern: "/search/test", ExpectedPattern: "/search/:name", OperateType: "search"},
		{Pattern: "/static/style.css", ExpectedPattern: "/static/*filepath", OperateType: "search"},
		{Pattern: "/insert", OperateType: "insert"},
		{Pattern: "/insert", ExpectedPattern: "/insert", OperateType: "search"},
	}

	for i, testCase := range testCases {
		if testCase.OperateType == "search" {
			parts := parsePattern(testCase.Pattern)
			node := root.search(parts, 0)
			if node == nil {
				t.Errorf("TestTrie #%v search node get error %v\n", i+1, errors.New("node is not exist"))
			} else {
				assert.Equal(t, testCase.ExpectedPattern, node.pattern)
			}
		}
		if testCase.OperateType == "insert" {
			parts := parsePattern(testCase.Pattern)
			root.insert(testCase.Pattern, parts, 0)
		}
	}
}
