// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
)

// Element defines the element of the CYK state table.
type Element struct {
	leftNonTerminal  string
	rightNonTerminal string

	begin int
	split int
	end   int
}

// NewBinary creates a binary element.
func NewBinary(l, r string) Element {
	return Element{leftNonTerminal: l, rightNonTerminal: r}
}

// NewUnary creates an unary element.
func NewUnary(l string) Element {
	return Element{leftNonTerminal: l}
}

// IsUnary returns true if the element is unary.
func (e Element) IsUnary() bool {
	return len(e.rightNonTerminal) == 0
}

// Set sets the limits for the item spans that are used by the CYK algorithm.
func (e Element) Set(begin, split, end int) Element {
	e.begin = begin
	e.split = split
	e.end = end
	return e
}

// CFG defines the Context-Free Grammar, which is specified by its production rules.
// It implements the Grammar interface.
type CFG struct {
	rules *Rules
}

// NewCFGrammar creates a new Context-Free Grammar. It loads the production rules from s.
func NewCFGrammar(s string) *CFG {
	return &CFG{rules: LoadRules(s)}
}

// BuildTrees computes the parse trees from the input items using the Lange-Leiss implementation
// of the CYK algorithm. The grammar is assumed to be in the binary normal form.
// Lange and Leiss, "To CNF or not to CNF? An Efficient Yet Presentable Version of the CYK Algorithm",
// Informatica Didactica 8 (2009).
func (g *CFG) BuildTrees(items Items) Trees {
	dim := items.Len()
	if dim == 0 {
		return nil
	}

	state := make([][]set.Set, dim)
	for i := 0; i < dim; i++ {
		state[i] = make([]set.Set, dim)
		for j := 0; j < dim; j++ {
			state[i][j] = set.New()
		}
	}
	children := make([][]map[string]Element, dim)
	for i := 0; i < dim; i++ {
		children[i] = make([]map[string]Element, dim)
		for j := 0; j < dim; j++ {
			children[i][j] = map[string]Element{}
		}
	}

	rules := g.rules

	k := 0
	for i := 0; i < dim; i++ {
		term := items[i].typ
		if rules.terminalRules[term].Empty() {
			continue
		}
		for A := range rules.terminalRules[term] {
			state[k][k].Add(A)
			children[k][k][A] = NewUnary(items[i].val).Set(k, k, k)
		}
		added := true
		for added {
			added = false
			for p, v := range rules.unaryRules {
				if A := p.leftNonTerminal; state[k][k][A] {
					for B := range v {
						if !state[k][k][B] {
							state[k][k].Add(B)
							children[k][k][B] = p.Set(k, k, k)
							added = true
						}
					}
				}
			}
		}
		k++
	}
	dim = k

	for span := 2; span <= dim; span++ {
		for begin := 0; begin <= dim-span; begin++ {
			end := begin + span - 1
			for split := begin; split < end; split++ {
				for B := range state[begin][split] {
					for C := range state[split+1][end] {
						p := NewBinary(B, C)
						for A := range rules.binaryRules[p] {
							state[begin][end].Add(A)
							children[begin][end][A] = p.Set(begin, split, end)
						}
					}
				}
			}
			foundUnaryDerivation := true
			for foundUnaryDerivation {
				foundUnaryDerivation = false
				for p, v := range rules.unaryRules {
					if A := p.leftNonTerminal; state[begin][end][A] {
						for B := range v {
							if !state[begin][end][B] {
								state[begin][end].Add(B)
								children[begin][end][B] = p.Set(begin, end, end)
								foundUnaryDerivation = true
							}
						}
					}
				}
			}
		}
	}

	// Build a parse tree stored in the children table:

	var iter func(n *Node, p Element)

	iter = func(n *Node, p Element) {
		node := NewNode(p.leftNonTerminal)
		n.left = node
		next, ok := children[p.begin][p.split][p.leftNonTerminal]
		if ok {
			iter(node, next)
		}

		if p.IsUnary() {
			return
		}

		node = NewNode(p.rightNonTerminal)
		n.right = node
		next, ok = children[p.split+1][p.end][p.rightNonTerminal]
		if ok {
			iter(node, next)
		}
	}

	trees := NewTrees()

	for k := 0; k < dim; k++ {
		score := 1.0 - 0.5*float64(k)/float64(dim)
		for i := 0; i <= k; i++ {
			if next, ok := children[i][dim+i-k-1]["S"]; ok {
				node := NewNode("S")
				iter(node, next)
				if node.Size() > 1 {
					tree := NewTree(node, score)
					trees = append(trees, tree)
				}
			}
		}
		if !trees.Empty() {
			break
		}
	}

	return trees
}
