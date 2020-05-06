// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNumber(t *testing.T) {
	a := assert.New(t)
	input := "12354"
	output := IsNumber(input)
	a.True(output)
}

func TestIsNotNumber(t *testing.T) {
	a := assert.New(t)
	input := "12t54"
	output := IsNumber(input)
	a.False(output)
}

func TestNormalizeText(t *testing.T) {
	a := assert.New(t)
	input := "Somewhere. - \"Here\" \n\nit's sunny.!"
	expected := "somewhere here its sunny"
	output := NormalizeText(input)
	a.Equal(expected, output)
}

func TestSplitWhitespace(t *testing.T) {
	a := assert.New(t)
	input := "somewhere here its sunny"
	expected := []string{"somewhere", "here", "its", "sunny"}
	output := SplitWhitespace(input)
	a.Equal(expected, output)
}

func TestSplitSentence(t *testing.T) {
	a := assert.New(t)
	input := "Somewhere here...  Its sunny! See you there."
	expected := []string{"Somewhere here", "Its sunny", "See you there", ""}
	output := SplitSentence(input)
	a.Equal(expected, output)
}

func TestCustomizeSlash(t *testing.T) {
	a := assert.New(t)
	input := "a /b"
	expected := []string{"a/b", "a / b", "a/ b", "a /b"}
	output := CustomizeSlash(input)
	a.Equal(expected, output)
}
