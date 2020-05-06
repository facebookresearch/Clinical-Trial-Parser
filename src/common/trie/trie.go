// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package trie

// Trie defines a search tree based on runes.
type Trie struct {
	root *Node
}

// Node defines a node in a trie.
type Node struct {
	value *Value
	nodes map[rune]*Node
	end   bool
}

// New creates a new trie.
func New() *Trie {
	root := &Node{nodes: make(map[rune]*Node), end: false}
	return &Trie{root}
}

// Put adds values to the trie.
func (t *Trie) Put(name string, vals ...string) {
	for _, val := range vals {
		v := NewValue(name, val)
		n := t.root
		for _, r := range v.val {
			m, ok := n.nodes[r]
			if !ok {
				m = &Node{nodes: make(map[rune]*Node), end: false}
				n.nodes[r] = m
			}
			n = m
		}
		n.value = v
		n.end = true
	}
}

// Get gets a value for the string s.
func (t *Trie) Get(s string) (*Value, bool) {
	n := t.root
	i := 0
	runes := []rune(s)
	for ; i < len(runes); i++ {
		m, ok := n.nodes['*']
		if ok {
			for ; i < len(runes) && rune(runes[i]) != ' '; i++ {
			}
			if i == len(runes) {
				break
			}
			i--
			n = m
		} else {
			r := runes[i]
			if m, ok := n.nodes[r]; ok {
				n = m
			} else {
				return nil, false
			}
		}
	}
	if !n.end {
		if m, ok := n.nodes['*']; ok {
			n = m
		}
	}
	return n.value, n.end
}

// Contains returns true if the string s is in the trie.
func (t *Trie) Contains(s string) bool {
	_, ok := t.Get(s)
	return ok
}

// Match returns true if the string matches a rune path in the trie regardless
// whether the trie contains a corresponding value.
func (t *Trie) Match(s string) bool {
	n := t.root
	i := 0
	runes := []rune(s)
	for ; i < len(runes); i++ {
		m, ok := n.nodes['*']
		if ok {
			for ; i < len(runes) && runes[i] != ' '; i++ {
			}
			if i == len(runes) {
				return true
			}
			i--
			n = m
		} else {
			r := runes[i]
			if m, ok := n.nodes[r]; ok {
				n = m
			} else {
				return false
			}
		}
	}

	if len(n.nodes) == 0 || n.end {
		return true
	}
	if _, ok := n.nodes[' ']; ok {
		return true
	}
	_, ok := n.nodes['*']
	return ok
}

// AllValues returns all the values of the node and its child nodes.
func (n *Node) AllValues() Values {
	results := Values{}
	if n.end == true {
		results = append(results, n.value)
	}
	for _, m := range n.nodes {
		prefixes := m.AllValues()
		results = append(results, prefixes...)
	}
	return results
}

// Autocomplete returns all the matching values of the prefix.
func (t *Trie) Autocomplete(prefix string) Values {
	n := t.root
	for _, r := range prefix {
		if m, ok := n.nodes[r]; ok {
			n = m
		} else {
			return Values{}
		}
	}
	values := n.AllValues()
	values.Sort()
	return values
}
