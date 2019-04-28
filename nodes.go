package main

import (
	"os"
)

// Node is an interface common for both known and unknown nodes
type Node interface {
	Path() string
}

// NodeUnknown is a path node which we cannot definitively know whether it exists, but
// we know its expected path
type NodeUnknown struct {
	path string
}

// Path satisfies the Node interface
func (n NodeUnknown) Path() string {
	return n.path
}

// NodeKnown is a path node which definitely exists
type NodeKnown struct {
	path   string
	st     os.FileInfo
	file   File
	target *ElaboratedPath
}

// Path satisfies the Node interface
func (n NodeKnown) Path() string {
	return n.path
}

// ElaboratedPath stores the nodes composing a path
type ElaboratedPath struct {
	parts  []Node
	groups map[string]bool
	owners map[string]bool
}

// GroupMaxLen gets the longest string length of the group name for each known node in the path
func (e *ElaboratedPath) GroupMaxLen() int {
	return maxLen(e.groups)
}

// OwnerMaxLen gets the longest string length of the owner name for each known node in the path
func (e *ElaboratedPath) OwnerMaxLen() int {
	return maxLen(e.owners)
}

func maxLen(m map[string]bool) int {
	maxLen := 0
	for item, _ := range m {
		if len(item) > maxLen {
			maxLen = len(item)
		}
	}
	return maxLen
}
