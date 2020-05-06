// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"fmt"
	"sort"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/relation"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/variables"
)

// Tree defines the parse tree.
type Tree struct {
	root  *Node
	score float64
}

// NewTree creates a new parse tree.
func NewTree(root *Node, score float64) *Tree {
	return &Tree{root: root, score: score}
}

// Size calculates the number of leafs (terminals).
func (t *Tree) Size() int {
	if t == nil || t.root == nil {
		return 0
	}
	return t.root.Size()
}

// Contains returns true if the tree t contains the tree v as a sub-tree or they are same.
func (t *Tree) Contains(v *Tree) bool {
	return t.root.Contains(v.root)
}

// String returns the string representation of the tree.
func (t *Tree) String() string {
	return fmt.Sprintf("{%q:%.3f,%q:%s}", "score", t.score, "tree", t.root.String())
}

// Relations converts the tree to 'or' and 'and' relations.
func (t *Tree) Relations() (relation.Relations, relation.Relations) {
	orRels, andRels := t.root.EvalRelations()
	orRels.SetScore(t.score)
	orRels.Sort()
	andRels.SetScore(t.score)
	andRels.Sort()
	return orRels, andRels
}

// Trees defines a slice of parse trees.
type Trees []*Tree

// NewTrees creates a new slice of trees.
func NewTrees() Trees {
	return make(Trees, 0)
}

// Dedupe deduplicates trees by removing exact duplicates and stumps
// that are contained by other trees.
func (ts *Trees) Dedupe() {
	a := *ts
	if len(a) < 2 {
		return
	}
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Size() > a[j].Size()
	})
	for i := len(a) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if a[j].Contains(a[i]) {
				a = append(a[:i], a[i+1:]...)
				break
			}
		}
	}
	*ts = a
}

// String returns the string representation of the trees.
func (ts Trees) String() string {
	switch len(ts) {
	case 0:
		return ""
	case 1:
		return ts[0].String()
	default:
		s := ts[0].String()
		for i := 1; i < len(ts); i++ {
			s += "\n" + ts[i].String()
		}
		return s
	}
}

// Relations converts the trees to 'or' and 'and' relations.
func (ts Trees) Relations() (relation.Relations, relation.Relations) {
	orRels := relation.NewRelations()
	andRels := relation.NewRelations()
	for _, t := range ts {
		or, and := t.Relations()
		orRels = append(orRels, or...)
		andRels = append(andRels, and...)
	}
	return orRels, andRels
}

// Empty tests whether ts has any trees in it.
func (ts Trees) Empty() bool {
	return len(ts) == 0
}

// Node defines a node in the parse tree, which is constructed by applying
// grammar production rules to the parsed items.
type Node struct {
	val   string
	left  *Node
	right *Node
}

// NewNode creates a new node.
func NewNode(val string) *Node {
	return &Node{val: val}
}

// Size calculates the number of leafs (terminals).
func (n *Node) Size() int {
	switch {
	case n == nil:
		return 0
	case n.left == nil && n.right == nil:
		return 1
	default:
		return n.left.Size() + n.right.Size()
	}
}

// Contains returns true if n contains m and its left and right children.
func (n *Node) Contains(m *Node) bool {
	switch {
	case m == nil:
		return true
	case n == nil:
		return false
	case n.val != m.val:
		return false
	default:
		return n.left.Contains(m.left) && n.right.Contains(m.right)
	}
}

// String returns the string representation of the node.
func (n *Node) String() string {
	s := fmt.Sprintf("{%q:%q", "value", n.val)
	if n.left != nil {
		s += fmt.Sprintf(",%q:%s", "left", n.left.String())
	}
	if n.right != nil {
		s += fmt.Sprintf(",%q:%s", "right", n.right.String())
	}
	s += "}"
	return s
}

// EvalVariable evaluates and returns the variable name stored in the terminal leaf.
func (n *Node) EvalVariable() string {
	variable := n.left.left.val
	if n.right != nil {
		variable += "/" + n.right.right.left.val
	}
	return variable
}

// EvalNums evaluates and returns the list of numbers stored in the terminal leafs.
func (n *Node) EvalNums() []string {
	set := set.New()
	var eval func(n *Node)
	eval = func(n *Node) {
		if n.left != nil {
			m := n.left
			if m.val == "N" {
				set.Add(m.left.val)
			} else {
				eval(m)
			}
		}
		if n.right != nil {
			m := n.right
			if m.val == "N" {
				set.Add(m.left.val)
			} else {
				eval(m)
			}
		}
	}
	eval(n)
	return set.Slice()
}

// EvalUnit evaluates and returns the unit stored in the terminal leaf.
func (n *Node) EvalUnit() string {
	unit := ""
	var eval func(n *Node)
	eval = func(n *Node) {
		if len(unit) > 0 {
			return
		}
		if n.left != nil {
			m := n.left
			if m.val == "U" {
				unit = m.left.val
				return
			}
			eval(m)
		}
		if n.right != nil {
			m := n.right
			if m.val == "U" {
				unit = m.left.val
				return
			}
			eval(m)
		}
	}
	eval(n)
	return unit
}

// EvalBound evaluates and returns the bound relation.
func (n *Node) EvalBound() (*relation.Limit, bool) {
	l := &relation.Limit{}
	lower := false

	if n.left.val == "L" && n.right.val == "T" {
		n.left, n.right = n.right, n.left
	}

	t := n.left
	if t.val != "T" {
		return nil, false
	}
	switch t.left.val {
	case "<":
		l.Value = "<"
		l.Incl = false
	case "≤":
		l.Value = "≤"
		l.Incl = true
	case ">":
		l.Value = ">"
		l.Incl = false
		lower = true
	case "≥":
		l.Value = "≥"
		l.Incl = true
		lower = true
	}
	num := n.right.left
	l.Value = num.left.val
	return l, lower
}

// EvalRange evaluates and returns the lower or upper bound of the range condition.
func (n *Node) EvalRange() *relation.Limit {
	l := &relation.Limit{}
	l.Incl = true
	nums := n.EvalNums()
	if len(nums) == 1 {
		l.Value = nums[0]
	}
	return l
}

// EvalRelation evaluates and returns the relation stored in the parse node based on the production rules.
func (n *Node) EvalRelation() (*relation.Relation, error) {
	r := relation.New()

	left := n.left
	right := n.right

	if left.val != "V" {
		left, right = right, left
	}
	if left.val != "V" {
		return nil, fmt.Errorf("bad or missing variable node: %s, %s", left.val, right.val)
	}

	r.Name = left.EvalVariable()
	r.ID, _ = variables.Get().ID(r.Name)

	// Check that the attribute node A exists:

	if right == nil || right.val != "A" {
		return r, nil
	}

	m := right

	r.Unit = m.EvalUnit()

	// Check the left branch of A:

	setBound := func(b *relation.Limit, lower bool) {
		switch {
		case lower:
			if r.Lower == nil {
				r.Lower = b
			}
		case !lower:
			if r.Upper == nil {
				r.Upper = b
			}
		}
	}

	left = m.left

	switch left.val {
	case "E":
		nums := m.EvalNums()
		r.Value = nums
		return r, nil
	case "B":
		b, lower := left.EvalBound()
		setBound(b, lower)
	case "L", "Y":
		r.Lower = left.EvalRange()
	}

	// Check the right branch of A:

	right = m.right

	if right == nil {
		return r, nil
	}

	switch right.val {
	case "W":
		b, lower := right.right.EvalBound()
		setBound(b, lower)
	case "B":
		b, lower := right.EvalBound()
		setBound(b, lower)
	case "Y":
		r.Upper = right.EvalRange()
	}

	return r, nil
}

// EvalRelations evaluates and returns the 'or' and 'and' relations stored in the parse node.
func (n *Node) EvalRelations() (relation.Relations, relation.Relations) {
	if n.left == nil {
		return relation.NewRelations(), relation.NewRelations()
	}
	switch {
	case n.left.val == "C" && n.right == nil:
		return n.left.EvalRelations()
	case n.left.val == "R" && n.right == nil:
		orRels := relation.NewRelations()
		andRels := relation.NewRelations()
		r, _ := n.left.EvalRelation()
		andRels = append(andRels, r)
		return orRels, andRels
	default:
		m := n.right
		conj := "and"
		if m.right != nil {
			conj = m.left.left.val
			m = m.right
		} else {
			m = m.left
		}
		orRels, andRels := n.left.EvalRelations()
		if r, err := m.EvalRelation(); err == nil {
			switch conj {
			case "or":
				orRels = append(orRels, r)
				orRels = append(orRels, andRels...)
				andRels = relation.NewRelations()
			default:
				andRels = append(andRels, r)
				andRels = append(andRels, orRels...)
				orRels = relation.NewRelations()
			}
		}
		return orRels, andRels
	}
}
