package gee

import (
	"fmt"
	"strings"
)

/*
	Our tri tree supports `:name(:lang etc.)` and `*filepath`.
	Tri tree has two important functions: Register and Matching.
	Register: if no node in current layer can match target part, insert a node in this layer
	Matching: search for every layers.
		Return if	(1)match *
					(2)match at len(parts)_th layer
					(3)fail to match
*/

type node struct {
	pattern  string  // registered pattern, only set at leaf nodes, e.g. /p/:lang
	part     string  // current node value, e.g. :lang
	children []*node // child nodes of current node, e.g. [doc, tutorial, intro]
	isWild   bool    // true if : or *
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// matchChild find 1st node that matches target. For inserting
func (n *node) matchChild(target string) *node {
	for _, child := range n.children {
		if child.part == target || child.isWild {
			return child
		}
	}
	return nil
}

// matChildren find all nodes that match target. For searching
func (n *node) matchChildren(target string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == target || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert() is called in router.go
func (n *node) insert(pattern string, parts []string, height int) {
	// only set pattern at leaf node, otherwise nil
	// when registered "/p/:lang/doc" and search for "/p/python".
	// Though "/p/python" can match with "/p/:lang" but pattern of ":/lang" is "". Hence, fail to match.
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	// such part is not found, insert a node
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'} // part[0] means the 1st character of current part
		n.children = append(n.children, child)
	}

	// insert recursively
	child.insert(pattern, parts, height+1)
}

// search() is called in router.go
func (n *node) search(parts []string, height int) *node {
	// Because length of filepath is not fixed, though we haven't traversed #height layers, search function also end and return.
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" { // if path until current layer didn't register
			return nil
		}
		return n
	}

	// search recursively
	part := parts[height] // value at current height
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// travel collects all registered pattern in tri tree
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
