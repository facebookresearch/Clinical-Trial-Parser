// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package set

import (
	"sort"
	"strings"
)

// Set implements a simple hash set for strings. The value is true,
// if the key exists. For a consistent behavior, do not explicitly
// set the value to false for a missing key.
type Set map[string]bool

func New(v ...string) Set {
	s := make(map[string]bool)
	for _, a := range v {
		s[a] = true
	}
	return s
}

func (s Set) Add(v ...string) {
	for _, a := range v {
		s[a] = true
	}
}

func (s Set) Remove(a string) bool {
	if !s[a] {
		return false
	}
	delete(s, a)
	return true
}

func (s Set) AddSet(other Set) {
	for a := range other {
		s[a] = true
	}
}

// Copy returns a copy of the set.
func (s Set) Copy() Set {
	copy := New()
	copy.AddSet(s)
	return copy
}

func (s Set) Get() (string, bool) {
	for k := range s {
		return k, true
	}
	return "", false
}

func (s Set) Contains(a string) bool {
	return s[a]
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Empty() bool {
	return s.Size() == 0
}

func (s Set) Intersection(other Set) int {
	a := s
	b := other
	if a.Size() > b.Size() {
		a, b = b, a
	}
	cnt := 0
	for e := range a {
		if b.Contains(e) {
			cnt++
		}
	}
	return cnt
}

func (s Set) Union(other Set) int {
	return s.Size() + other.Size() - s.Intersection(other)
}

func (s Set) Jaccard(other Set) float64 {
	intersection := s.Intersection(other)
	if intersection == 0 {
		return 0
	}
	union := s.Union(other)
	return float64(intersection) / float64(union)
}

func (s Set) Slice() []string {
	list := make([]string, 0, s.Size())
	for key, _ := range s {
		list = append(list, key)
	}
	sort.Strings(list)
	return list
}

func (s Set) String() string {
	list := s.Slice()
	return strings.Join(list, ", ")
}
