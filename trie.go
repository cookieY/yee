package knocker

import (
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
	priority int8
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.pattern == part || child.isWild {
			if len(n.part) > 0 {
				if n.part[0] != ':' {
					return child
				}
			} else {
				return child
			}
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
		if part[0] == ':' || part[0] == '*' {
			child = &node{part: part, isWild: true, priority: 2}
		} else {
			child = &node{part: part, isWild: false, priority: 1}
		}
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

	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}

func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
