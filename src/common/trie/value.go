// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package trie

import (
	"sort"
)

// Value defines the value stored in the trie node.
type Value struct {
	name string
	val  string
}

// NewValue creates a new value object.
func NewValue(name string, val string) *Value {
	return &Value{name: name, val: val}
}

// Name returns the name of the value object.
func (v Value) Name() string {
	return v.name
}

// Values defines the slice of values.
type Values []*Value

// Sort sorts values in ascending order by name.
func (vs Values) Sort() {
	sort.SliceStable(vs, func(i, j int) bool {
		return vs[i].name < vs[j].name
	})
}
