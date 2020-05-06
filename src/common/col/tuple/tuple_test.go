// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package tuple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortTuple(t *testing.T) {
	a := assert.New(t)

	expected := New("a", "a", "ac", "d")
	actual := New("d", "a", "ac", "a")
	actual.Sort()

	a.Equal(expected, actual)
}

func TestSortTuples(t *testing.T) {
	a := assert.New(t)

	expected := Tuples{
		New("Canada", "ON"),
		New("United States", "CA"),
		New("United States", "NY"),
		New("Xanada", "XO"),
	}
	actual := Tuples{
		expected[3],
		expected[2],
		expected[1],
		expected[0],
	}
	actual.Sort()

	a.Equal(expected, actual)
}
