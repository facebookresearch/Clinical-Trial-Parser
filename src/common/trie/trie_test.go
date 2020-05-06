// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var dictionary *Trie

func init() {
	dictionary = New()
	dictionary.Put("A", "abcd")
	dictionary.Put("B", "abd")
	dictionary.Put("C", "x")
	dictionary.Put("D", "a1234* 789")
	dictionary.Put("E", "1234*")
	dictionary.Put("F", "wx y*")
	dictionary.Put("kg/m2", "kg/m²")
}

func TestContains(t *testing.T) {
	a := assert.New(t)

	input := "abcd"
	actual := dictionary.Contains(input)
	a.True(actual)
}

func TestDoesntContain(t *testing.T) {
	a := assert.New(t)

	input := "abcde"
	actual := dictionary.Contains(input)
	a.False(actual)
}

func TestSpecialUnicodeContains(t *testing.T) {
	a := assert.New(t)

	input := "kg/m²"
	actual := dictionary.Contains(input)
	a.True(actual)
}

func TestMatch(t *testing.T) {
	a := assert.New(t)

	input := "wx"
	actual := dictionary.Match(input)
	a.True(actual)

	input = "wx yz"
	actual = dictionary.Match(input)
	a.True(actual)
}

func TestNoMatch(t *testing.T) {
	a := assert.New(t)

	input := "a"
	actual := dictionary.Match(input)
	a.False(actual)
}

func TestWild(t *testing.T) {
	a := assert.New(t)

	input := "123"
	_, ok := dictionary.Get(input)
	a.False(ok)

	input = "1234"
	actual, ok := dictionary.Get(input)
	a.True(ok)
	a.Equal("1234*", actual.val)

	input = "12345"
	actual, ok = dictionary.Get(input)
	a.True(ok)
	a.Equal("1234*", actual.val)

	input = "a123456 789"
	actual, ok = dictionary.Get(input)
	a.True(ok)
	a.Equal("a1234* 789", actual.val)
}

func TestMatchWild(t *testing.T) {
	a := assert.New(t)

	input := "123"
	actual := dictionary.Match(input)
	a.False(actual)

	input = "1234"
	actual = dictionary.Match(input)
	a.True(actual)

	input = "12345"
	actual = dictionary.Match(input)
	a.True(actual)

	input = "a123456 789"
	actual = dictionary.Match(input)
	a.True(actual)
}

func TestAutocomplete(t *testing.T) {
	a := assert.New(t)

	input := "ab"
	expected := Values{NewValue("A", "abcd"), NewValue("B", "abd")}
	actual := dictionary.Autocomplete(input)
	a.Equal(expected, actual)
}
