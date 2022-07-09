package web

import (
	"log"
	"strings"
)

type node struct {
	pattern  string  // router path to be matched
	part     string  // matched part in router path
	children []*node // sub path of mathch router
	isWild   bool    // match wildcard */: true
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	matched := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			log.Printf("matchChildren(%q): child: %v", part, child)
			matched = append(matched, child)
		}
	}
	return matched
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.part = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: string(part[0]) == ":" || string(part[0]) == "*",
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
	log.Printf("search part: %v", part)

	children := n.matchChildren(part)
	for _, child := range children {
		log.Printf("search child: %v", child)
		result := child.search(parts, height+1)
		log.Printf("search result: %v", result)
		if result != nil {
			return result
		}
	}
	return nil
}
