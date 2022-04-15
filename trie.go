package gem

import (
	"strings"
)

type methodTree struct {
	method string
	root   *node
}

type methodTrees []methodTree

func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
	handlers []HandlerFunc
}

func (n *node) insert(pattern string, parts []string, height int, handlers []HandlerFunc) {
	if len(parts) == height {
		n.pattern = pattern
		n.handlers = handlers
		return
	}
	part := parts[height]
	var child *node
	for _, c := range n.children {
		if c.part == part || c.isWild {
			child = c
			break
		}
	}
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1, handlers)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	for _, child := range children {
		node := child.search(parts, height+1)
		if node != nil {
			return node
		}
	}
	return nil
}

func (n *node) getValue(path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	node := n.search(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func parsePattern(pattern string) []string {
	values := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range values {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
