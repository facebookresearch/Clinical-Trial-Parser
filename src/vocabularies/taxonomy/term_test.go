// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package taxonomy

import (
	"testing"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"

	"github.com/stretchr/testify/assert"
)

var (
	ATerm  Term
	BTerm  Term
	CTerm  Term
	MTerm  Term
	ACTerm Term
)

func initTerms() {
	ATerm = NewTerm("a", 1, set.New("a"), set.New("A"))
	BTerm = NewTerm("b", 1, set.New("b"), set.New("B"))
	CTerm = NewTerm("c", 1, set.New("c"), set.New("C"))
	MTerm = NewTerm("a", 1, set.New("c"), set.New("C"))
	ACTerm = NewTerm("a", 1, set.New("a", "c"), set.New("A", "C"))
}

func TestDistinctDedupe(t *testing.T) {
	initTerms()
	a := assert.New(t)

	input := Terms{ATerm, BTerm, CTerm}
	expected := Terms{ATerm, BTerm, CTerm}
	actual := input.Dedupe()

	a.Equal(expected, actual)
}

func TestMergedDedupe(t *testing.T) {
	initTerms()
	a := assert.New(t)

	input := Terms{ATerm, MTerm, BTerm}
	expected := Terms{ACTerm, BTerm}
	actual := input.Dedupe()

	a.Equal(expected, actual)
}

func TestMixedDedupe(t *testing.T) {
	initTerms()
	a := assert.New(t)

	input := Terms{BTerm, ATerm, MTerm}
	expected := Terms{BTerm, ACTerm}
	actual := input.Dedupe()

	a.Equal(expected, actual)
}

func TestDuplicateDedupe(t *testing.T) {
	initTerms()
	a := assert.New(t)

	input := Terms{ATerm, MTerm, MTerm}
	expected := Terms{ACTerm}
	actual := input.Dedupe()

	a.Equal(expected, actual)
}

func TestPassFilter(t *testing.T) {
	initTerms()
	a := assert.New(t)

	input := Terms{ATerm, BTerm, CTerm}
	expected := Terms{CTerm}
	validCategories := set.New("c")
	actual := input.PassFilter(validCategories)

	a.Equal(expected, actual)
}
