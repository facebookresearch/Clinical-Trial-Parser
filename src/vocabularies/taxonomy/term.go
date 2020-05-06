// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package taxonomy

import (
	"fmt"
	"math"
	"sort"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/text"
)

// Term defines a vocabulary term or concept with a Key (e.g., name)
// and Value (e.g., score).
type Term struct {
	Key         string
	Normalized  string
	Value       float64
	Categories  set.Set
	TreeNumbers set.Set
}

// NewTerm creates a term.
func NewTerm(k string, v float64, c set.Set, t set.Set) Term {
	return Term{Key: k, Value: v, Categories: c, TreeNumbers: t}
}

// String returns a string representation of the term.
func (t Term) String() string {
	return fmt.Sprintf("%s: %.2f | %s | %s", t.Key, t.Value, t.Categories.String(), t.TreeNumbers.String())
}

// PassFilter returns true if any of the terms categories are in filter.
func (t Term) PassFilter(filter set.Set) bool {
	if filter.Empty() {
		return true
	}
	if t.Categories.Intersection(filter) > 0 {
		return true
	}
	return false
}

// TrimCategories removes categories and tree numbers that are not in the categories set.
func (t Term) TrimCategories(categories set.Set) Term {
	if categories.Empty() {
		return t
	}
	for c := range t.Categories {
		if !categories[c] {
			delete(t.Categories, c)
		}
	}
	for tn := range t.TreeNumbers {
		if c := text.LetterPrefix(tn); !categories[c] {
			delete(t.TreeNumbers, tn)
		}
	}
	return t
}

// Terms define a slice of terms.
type Terms []Term

// Default creates a slice of terms with a 0 value term.
func Default(k, n string) Terms {
	return Terms{Term{Key: k, Normalized: n, Categories: set.New(), TreeNumbers: set.New()}}
}

//NewTerms creates a slice of terms.
func NewTerms(cap int) Terms {
	return make(Terms, cap)
}

// Len returns the number of terms.
func (ts Terms) Len() int {
	return len(ts)
}

// MaxValue returns the largest value. Terms are assumed to be
// sorted in reverse order by Value.
func (ts Terms) MaxValue() float64 {
	if ts.Len() > 0 {
		return ts[0].Value
	}
	return -1
}

// MaxKey returns the key of the term with the largest value.
// Terms are assumed to be sorted in reverse order by Value.
func (ts Terms) MaxKey() string {
	if ts.Len() > 0 {
		return ts[0].Key
	}
	return ""
}

// Normalized returns the normalized key.
func (ts Terms) Normalized() string {
	if ts.Len() > 0 {
		return ts[0].Normalized
	}
	return ""
}

// Keys returns the the unique keys of the terms.
func (ts Terms) Keys() []string {
	keys := set.New()
	for _, t := range ts {
		keys.Add(t.Key)
	}
	return keys.Slice()
}

// Categories returns the unique categories of the terms.
func (ts Terms) Categories() []string {
	cat := set.New()
	for _, t := range ts {
		cat.AddSet(t.Categories)
	}
	return cat.Slice()
}

// TreeNumbers returns the unique tree numbers of the terms.
func (ts Terms) TreeNumbers() []string {
	addr := set.New()
	for _, t := range ts {
		addr.AddSet(t.TreeNumbers)
	}
	return addr.Slice()
}

// String returns a string representation of the terms.
func (ts Terms) String() string {
	switch len(ts) {
	case 0:
		return ""
	case 1:
		return ts[0].String()
	default:
		s := ts[0].String()
		for i := 1; i < ts.Len(); i++ {
			s += fmt.Sprintf("\n%s", ts[i].String())
		}
		return s
	}
}

// Dedupe de-duplicates terms by their keys. Categories and tree numbers
// of duplicate terms are joined.
func (ts Terms) Dedupe() Terms {
	if ts.Len() < 2 {
		return ts
	}
	n := 0
	for i := 1; i < ts.Len(); i++ {
		if ts[n].Key != ts[i].Key {
			n++
			ts[n] = ts[i]
		} else {
			ts[n].Value = math.Max(ts[n].Value, ts[i].Value)
			ts[n].Categories.AddSet(ts[i].Categories)
			ts[n].TreeNumbers.AddSet(ts[i].TreeNumbers)
		}
	}
	return ts[:n+1]
}

// PassFilter keeps terms whose categories are in categories filter.
func (ts Terms) PassFilter(categories set.Set) Terms {
	if categories.Empty() {
		return ts
	}
	n := 0
	for _, t := range ts {
		if t.Categories.Intersection(categories) > 0 {
			ts[n] = t.TrimCategories(categories)
			n++
		}
	}
	return ts[:n]
}

// TopDelta keep terms whose terms are within d of the top term.
func (ts Terms) TopDelta(d float64) Terms {
	if ts.Len() < 2 {
		return ts
	}
	i := 1
	for ; i < ts.Len(); i++ {
		if ts[0].Value-ts[i].Value > d {
			break
		}
	}
	return ts[:i]
}

// SortByValue sorts terms by value in reverse order.
func (ts Terms) SortByValue() Terms {
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Value > ts[j].Value
	})
	return ts
}

// SortByKey sorts terms by key.
func (ts Terms) SortByKey() Terms {
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Key < ts[j].Key
	})
	return ts
}
