// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package lsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var minHash MinHash

func init() {
	minHash = New(3, 16)
}

func TestSameness(t *testing.T) {
	a := assert.New(t)
	s1 := "autism spectrum disorder"
	s2 := "autism spectrum disorder"
	a.True(minHash.IsSimilar(s1, s2))
}

func TestSimilarity(t *testing.T) {
	a := assert.New(t)
	s1 := "autism spectrum disorder"
	s2 := "autism spectrum"
	a.True(minHash.IsSimilar(s1, s2))
}

func TestDifference(t *testing.T) {
	a := assert.New(t)
	s1 := "autism spectrum disorder"
	s2 := "gastrointestinal disorders"
	a.False(minHash.IsSimilar(s1, s2))
}
