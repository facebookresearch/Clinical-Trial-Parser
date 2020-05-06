// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package taxonomy

import (
	"fmt"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/lsh"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/fio"

	"github.com/golang/glog"
)

const (
	capacity = 20
	buffSize = 10000
	minScore = 0.4
)

// Taxonomy defines a taxonomy for a vocabulary.
type Taxonomy struct {
	root      *Node
	normalize Normalizer

	baseIndex []int
	hashIndex map[string][]int
	minHash   lsh.MinHash

	capacity int
	buffSize int
	minScore float64
}

// New creates a new taxonomy.
func New(r *Node) *Taxonomy {
	return &Taxonomy{root: r, normalize: identity, capacity: capacity, buffSize: buffSize, minScore: minScore}
}

// SetQueueCapacity sets the capacity of the search priority queue.
func (t *Taxonomy) SetQueueCapacity(c int) {
	t.capacity = c
}

// SetBuffSize sets the buffer size of the search channel.
func (t *Taxonomy) SetBuffSize(b int) {
	t.buffSize = b
}

// SetMinScore sets the minimum score below which terms are disregarded.
func (t *Taxonomy) SetMinScore(p float64) {
	t.minScore = p
}

// AddNodes adds nodes to the taxonomy. Nodes with the same name are joined.
func (t *Taxonomy) AddNodes(ns Nodes) int {
	cnt := 0
	for _, n := range ns {
		if t.AddNode(n) {
			cnt++
		}
	}
	return cnt
}

// AddNode adds a node to the taxonomy. Nodes with the same name are joined.
func (t *Taxonomy) AddNode(n *Node) bool {
	if ok := t.root.Update(n); !ok {
		t.root.AddChild(n)
		return true
	}
	return false
}

// Normalize normalizes the node synonyms.
func (t *Taxonomy) Normalize(f Normalizer) {
	t.normalize = f
	t.root.normalize(f)
}

// Info prints basic information about the taxonomy.
func (t *Taxonomy) Info() {
	fmt.Printf("Descriptors: %6d\n", t.root.Size(0, 1))
	fmt.Printf("Concepts:    %6d\n", t.root.Size(0, 2))
	fmt.Printf("Terms:       %6d\n", t.root.Leafs())
}

// Store stores the taxonomy to a file.
func (t *Taxonomy) Store(fname, sep string) error {
	writer := fio.Writer(fname)
	defer writer.Close()

	for _, n := range t.root.children {
		for _, m := range n.children {
			syn := strings.Join(m.synonyms.Slice(), "|")
			for tn := range m.treeNumbers {
				if _, err := fmt.Fprintf(writer, "%s%s%s%s%s%s%s\n", n.name, sep, m.name, sep, tn, sep, syn); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Candidates finds a set of candidate nodes for given string.
func (t *Taxonomy) Candidates(s string) []int {
	candidates := make(map[int]bool)
	codes := t.minHash.HashCodes(s)
	for code := range codes {
		es := t.hashIndex[code]
		for _, e := range es {
			candidates[e] = true
		}
	}
	list := make([]int, 0, len(candidates))
	for e := range candidates {
		list = append(list, e)
	}
	return list
}

// SetBaseIndex sets the index to a slice.
func (t *Taxonomy) SetBaseIndex() {
	baseIndex := make([]int, t.root.children.Len())
	for i := range baseIndex {
		baseIndex[i] = i
	}
	t.baseIndex = baseIndex
	t.minHash = lsh.New(3, 16) // For computing similarity scores.
	t.hashIndex = nil
}

// SetHashIndex sets the indexing to LSH.
func (t *Taxonomy) SetHashIndex(rows, bands int) {
	t.SetBaseIndex()

	fmt.Print("Indexing ... ")
	minHash := lsh.New(rows, bands)
	hashIndex := make(map[string][]int)
	for i, n := range t.root.children {
		codes := n.hashCodes(minHash)
		for c := range codes {
			hashIndex[c] = append(hashIndex[c], i)
		}
	}
	t.hashIndex = hashIndex
	t.minHash = minHash
	fmt.Println("indexed")
}

// getMatchIndices gets the indices of the candidate nodes.
func (t *Taxonomy) getMatchIndices(s string) []int {
	if len(t.hashIndex) == 0 {
		return t.baseIndex
	}
	if indices := t.Candidates(s); len(indices) > 0 {
		return indices
	}
	return t.baseIndex
}

// Match matches a string to terms in the taxonomy.
func (t *Taxonomy) Match(s string, d float64, filter set.Set) Terms {
	if len(t.baseIndex) == 0 && len(t.hashIndex) == 0 {
		glog.Fatal("Search index not set.")
	}
	priority := NewPriority(t.capacity)
	q := make(chan Term, t.buffSize)

	nsorted, n := t.normalize(s)

	if len(nsorted) == 0 {
		return Default(s, n)
	}
	indices := t.getMatchIndices(nsorted)

	go t.root.walk(nsorted, indices, q, t.minHash, t.minScore)

	for p := range q {
		if p.PassFilter(filter) {
			priority.Insert(p.TrimCategories(filter))
		}
	}

	terms := priority.Terms().SortByKey().Dedupe().SortByValue().TopDelta(d)
	if terms.Len() > 0 {
		terms[0].Normalized = n
	} else {
		terms = Default(s, n)
	}
	return terms
}

// MatchNode matches a string to terms in the taxonomy.
func (t *Taxonomy) MatchNode(n *Node, d float64, filter set.Set) Terms {
	if len(t.baseIndex) == 0 || len(t.hashIndex) == 0 {
		glog.Fatal("Search index not set.")
	}

	empty := set.New()
	matches := t.Match(n.name, d, empty)
	for s := range n.synonyms {
		matches = append(matches, t.Match(s, d, filter)...)
	}
	return matches.SortByKey().Dedupe().SortByValue()
}
