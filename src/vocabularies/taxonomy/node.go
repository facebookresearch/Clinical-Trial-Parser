// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package taxonomy

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/lsh"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/slice"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"

	"github.com/golang/glog"
)

// Node defines a node in a taxonomy.
type Node struct {
	name        string
	children    Nodes
	synonyms    set.Set
	treeNumbers set.Set
}

// NewNode creates a new node.
func NewNode(s string) *Node {
	return &Node{name: s, synonyms: set.New(), treeNumbers: set.New()}
}

// Name returns the node's name.
func (n *Node) Name() string {
	return n.name
}

// Synonyms returns all synonyms of the node and its child nodes.
func (n *Node) Synonyms() set.Set {
	set := set.New()
	set.AddSet(n.synonyms)
	for _, m := range n.children {
		set.AddSet(m.Synonyms())
	}
	return set
}

// TreeNumbers returns all tree numbers of the node and its child nodes.
func (n *Node) TreeNumbers() set.Set {
	set := set.New()
	set.AddSet(n.treeNumbers)
	for _, m := range n.children {
		set.AddSet(m.treeNumbers)
	}
	return set
}

// Categories returns all categories of the node and its child nodes.
func (n *Node) Categories() set.Set {
	categories := set.New()
	for _, m := range n.children {
		categories.AddSet(m.Categories())
	}
	for a := range n.treeNumbers {
		categories.Add(text.LetterPrefix(a))
	}
	return categories
}

// Size is the number of nodes at level l.
func (n *Node) Size(i, l int) int {
	if n.children == nil {
		return 1
	}
	if i == l {
		return 1
	}
	leafs := 0
	for _, m := range n.children {
		leafs += m.Size(i+1, l)
	}
	return leafs
}

// Leafs is the number of synonyms of leaf nodes.
func (n *Node) Leafs() int {
	if n.children == nil {
		return len(n.synonyms) - 1
	}
	leafs := 0
	for _, m := range n.children {
		leafs += m.Leafs()
	}
	return leafs
}

// AddChild adds a node to the nodes children.
func (n *Node) AddChild(c *Node) {
	n.children = append(n.children, c)
}

// AddSynonym add a slice of synonyms to the node.
func (n *Node) AddSynonym(ss ...string) {
	n.synonyms.Add(ss...)
}

// AddTreeNumber add a slice of tree numbers to the node.
func (n *Node) AddTreeNumber(tn ...string) {
	n.treeNumbers.Add(tn...)
}

func equals(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

// Update finds and updates a node. If no matching node is found, false is returned.
func (n *Node) Update(m *Node) bool {
	if n.children == nil && equals(n.name, m.name) {
		n.synonyms.AddSet(m.synonyms)
		n.treeNumbers.AddSet(m.treeNumbers)
		return true
	}
	status := false
	for _, c := range n.children {
		status = status || c.Update(m)
	}
	return status
}

// normalize normalizes the node's and its child nodes' synonyms.
func (n *Node) normalize(f Normalizer) {
	synonyms := set.New()
	for s := range n.synonyms {
		synonyms.Add(f(s))
	}
	n.synonyms = synonyms

	for i := range n.children {
		n.children[i].normalize(f)
	}
}

// walk walks the node and its child nodes and returns terms with match score.
func (n *Node) walk(s string, indices []int, q chan<- Term, minHash lsh.MinHash, minScore float64) {
	wait := &sync.WaitGroup{}
	for _, i := range indices {
		wait.Add(1)
		m := n.children[i]
		go func(m *Node) {
			defer wait.Done()
			m.match(s, q, minHash, minScore)
		}(m)
	}

	go func() {
		wait.Wait()
		close(q)
	}()
}

// match returns the node and its child nodes with the match scores.
func (n *Node) match(s string, q chan<- Term, minHash lsh.MinHash, minScore float64) {
	for syn := range n.synonyms {
		if score := minHash.Similarity(s, syn); score >= minScore {
			q <- NewTerm(n.name, score, n.Categories(), n.TreeNumbers().Copy())
		}
	}
	for _, m := range n.children {
		m.match(s, q, minHash, minScore)
	}
}

// hashCodes returns the has codes of the node's synonyms.
func (n *Node) hashCodes(h lsh.MinHash) set.Set {
	codes := set.New()
	for s := range n.synonyms {
		codes.AddSet(h.HashCodes(s))
	}
	return codes
}

// Nodes defines a slice of nodes.
type Nodes []*Node

// NewNodes creates a slice of nodes.
func NewNodes() Nodes {
	return Nodes{}
}

// Len returns the number of nodes.
func (ns Nodes) Len() int {
	return len(ns)
}

// LoadNodes loads nodes from a file.
func LoadNodes(fnames ...string) Nodes {
	index := make(map[string]int)
	nodes := NewNodes()

	for _, fname := range fnames {
		file, err := os.Open(fname)
		if err != nil {
			glog.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineCnt := 0

		for scanner.Scan() {
			lineCnt++
			line := scanner.Text()
			line = strings.TrimSpace(line)

			if len(line) > 0 {
				values := strings.Split(line, "\t")
				slice.TrimSpace(values)
				values = slice.RemoveEmpty(values)
				var treeNumbers []string
				switch len(values) {
				case 0, 1:
					glog.Fatalf("%s: Too few columns: line %d: '%s'\n", fname, lineCnt, line)
				case 2:
					treeNumbers = nil
				case 3:
					treeNumbers = strings.Split(values[2], "|")
				default:
					glog.Fatalf("%s: Too many columns: line %d: '%s'\n", fname, lineCnt, line)
				}
				conceptName := values[0]
				lowercase := strings.ToLower(conceptName)
				synonym := values[1]

				if i, ok := index[lowercase]; ok {
					nodes[i].AddSynonym(synonym)
					nodes[i].AddTreeNumber(treeNumbers...)
				} else {
					node := NewNode(conceptName)
					node.AddSynonym(conceptName)
					node.AddSynonym(synonym)
					node.AddTreeNumber(treeNumbers...)
					index[lowercase] = nodes.Len()
					nodes = append(nodes, node)
				}
			}
		}
	}

	return nodes
}
