// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	a := assert.New(t)
	s := New("a", "b")
	a.True(s.Contains("a"))
}

func TestNotContain(t *testing.T) {
	a := assert.New(t)
	s := New("a")
	a.False(s.Contains("b"))
}

func TestRemove(t *testing.T) {
	a := assert.New(t)
	s := New("a", "c")
	a.False(s.Remove("b"))
	a.True(s.Remove("a"))
	a.Equal(s, New("c"))
}

func TestCopy(t *testing.T) {
	a := assert.New(t)
	expected := New("1", "2")
	actual := expected.Copy()
	a.Equal(expected, actual)
}

func TestIntersection(t *testing.T) {
	a := assert.New(t)
	s1 := New("a", "b")
	s2 := New("b", "c")
	a.Equal(s1.Intersection(s2), 1)
}

func TestUnion(t *testing.T) {
	a := assert.New(t)
	s1 := New("a", "b")
	s2 := New("b", "c")
	a.Equal(s1.Union(s2), 3)
}
