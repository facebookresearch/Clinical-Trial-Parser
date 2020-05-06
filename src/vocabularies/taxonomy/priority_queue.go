// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package taxonomy

// Priority defines a priority queue.
type Priority struct {
	capacity int
	size     int
	terms    Terms
}

// NewPriority creates a new priority queue.
func NewPriority(cap int) *Priority {
	return &Priority{capacity: cap, size: 0, terms: NewTerms(cap + 1)}
}

// Size returns the number of terms in the queue.
func (p *Priority) Size() int {
	return p.size
}

// Terms returns the terms in the queue.
func (p *Priority) Terms() Terms {
	return p.terms[:p.size]
}

// Insert inserts a term to the queue.
func (p *Priority) Insert(t Term) {
	j := p.size
	for j > 0 {
		if t.Value > p.terms[j-1].Value {
			p.terms[j] = p.terms[j-1]
			j--
		} else {
			break
		}
	}
	p.terms[j] = t
	if p.size < p.capacity {
		p.size++
	}
}
