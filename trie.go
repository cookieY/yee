package yee

import (
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.pattern == part {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	node := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			node = append(node, child)
		}
	}
	return node
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {

	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]

	children := n.matchChildren(part)

	var nodes []*node

	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			nodes = append(nodes, res)
		}
	}

	index := 0

	for _, i := range nodes {

		if i.isWild {
			index++
		} else {
			return i
		}
	}

	if index > 0 && nodes != nil {
		return nodes[index-1]
	}
	return nil
}
