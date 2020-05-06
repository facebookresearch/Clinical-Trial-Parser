// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestA1CLexer(t *testing.T) {
	a := assert.New(t)

	input := "a1c (hemoglobin a1c) greater than or equal to 5.0%."
	expected := Tokens{
		NewToken(tokenIdentifier, 0, "a1c"),
		NewToken(tokenLeftParenthesis, 4, "("),
		NewToken(tokenIdentifier, 5, "hemoglobin"),
		NewToken(tokenIdentifier, 16, "a1c"),
		NewToken(tokenRightParenthesis, 19, ")"),
		NewToken(tokenGreaterComparison, 21, "greater"),
		NewToken(tokenComparison, 29, "than"),
		NewToken(tokenConjunction, 34, "or"),
		NewToken(tokenIdentifier, 37, "equal"),
		NewToken(tokenIdentifier, 43, "to"),
		NewToken(tokenNumber, 46, "5.0"),
		NewToken(tokenUnit, 49, "%"),
		NewToken(tokenPunctuation, 50, "."),
	}
	actual := NewLexer(input).Drain()
	a.Equal(expected, actual)
}

func TestWeightLexer(t *testing.T) {
	a := assert.New(t)

	input := "patient > 50 kg in weight"
	expected := Tokens{
		NewToken(tokenIdentifier, 0, "patient"),
		NewToken(tokenComparison, 8, ">"),
		NewToken(tokenNumber, 10, "50"),
		NewToken(tokenIdentifier, 13, "kg"),
		NewToken(tokenIdentifier, 16, "in"),
		NewToken(tokenIdentifier, 19, "weight"),
	}
	actual := NewLexer(input).Drain()
	a.Equal(expected, actual)
}

func TestBMILexer(t *testing.T) {
	a := assert.New(t)

	input := "bmi > 50 kg/m²."
	expected := Tokens{
		NewToken(tokenIdentifier, 0, "bmi"),
		NewToken(tokenComparison, 4, ">"),
		NewToken(tokenNumber, 6, "50"),
		NewToken(tokenIdentifier, 9, "kg/m²"),
		NewToken(tokenPunctuation, 15, "."),
	}
	actual := NewLexer(input).Drain()
	a.Equal(expected, actual)
}

func TestBPLexer(t *testing.T) {
	a := assert.New(t)

	input := "uncontrolled hypertension (systolic blood pressure (sbp) >140mmhg, diastolic blood pressure (dbp) >90mmhg)"
	expected := Tokens{
		NewToken(tokenIdentifier, 0, "uncontrolled"),
		NewToken(tokenIdentifier, 13, "hypertension"),
		NewToken(tokenLeftParenthesis, 26, "("),
		NewToken(tokenIdentifier, 27, "systolic"),
		NewToken(tokenIdentifier, 36, "blood"),
		NewToken(tokenIdentifier, 42, "pressure"),
		NewToken(tokenLeftParenthesis, 51, "("),
		NewToken(tokenIdentifier, 52, "sbp"),
		NewToken(tokenRightParenthesis, 55, ")"),
		NewToken(tokenComparison, 57, ">"),
		NewToken(tokenNumber, 58, "140"),
		NewToken(tokenIdentifier, 61, "mmhg"),
		NewToken(tokenChar, 65, ","),
		NewToken(tokenIdentifier, 67, "diastolic"),
		NewToken(tokenIdentifier, 77, "blood"),
		NewToken(tokenIdentifier, 83, "pressure"),
		NewToken(tokenLeftParenthesis, 92, "("),
		NewToken(tokenIdentifier, 93, "dbp"),
		NewToken(tokenRightParenthesis, 96, ")"),
		NewToken(tokenComparison, 98, ">"),
		NewToken(tokenNumber, 99, "90"),
		NewToken(tokenIdentifier, 101, "mmhg"),
		NewToken(tokenRightParenthesis, 105, ")"),
	}
	actual := NewLexer(input).Drain()
	a.Equal(expected, actual)
}

func TestASTLexer(t *testing.T) {
	a := assert.New(t)

	input := "aspartate aminotransferase (ast)/alanine aminotransferase (alt) ≤ 2.0 x upper limits of normal"
	expected := Tokens{
		NewToken(tokenIdentifier, 0, "aspartate"),
		NewToken(tokenIdentifier, 10, "aminotransferase"),
		NewToken(tokenLeftParenthesis, 27, "("),
		NewToken(tokenIdentifier, 28, "ast"),
		NewToken(tokenRightParenthesis, 31, ")"),
		NewToken(tokenSlash, 32, "/"),
		NewToken(tokenIdentifier, 33, "alanine"),
		NewToken(tokenIdentifier, 41, "aminotransferase"),
		NewToken(tokenLeftParenthesis, 58, "("),
		NewToken(tokenIdentifier, 59, "alt"),
		NewToken(tokenRightParenthesis, 62, ")"),
		NewToken(tokenComparison, 64, "≤"),
		NewToken(tokenNumber, 68, "2.0"),
		NewToken(tokenIdentifier, 72, "x"),
		NewToken(tokenIdentifier, 74, "upper"),
		NewToken(tokenIdentifier, 80, "limits"),
		NewToken(tokenIdentifier, 87, "of"),
		NewToken(tokenIdentifier, 90, "normal"),
	}
	actual := NewLexer(input).Drain()
	a.Equal(expected, actual)
}

func TestWBCLexer(t *testing.T) {
	a := assert.New(t)

	input := "wbc > 3,000/ul"
	expected := Tokens{
		NewToken(tokenIdentifier, 0, "wbc"),
		NewToken(tokenComparison, 4, ">"),
		NewToken(tokenNumber, 6, "3,000"),
		NewToken(tokenSlash, 11, "/"),
		NewToken(tokenIdentifier, 12, "ul"),
	}
	actual := NewLexer(input).Drain()
	a.Equal(expected, actual)
}

func TestNumberLexer(t *testing.T) {
	a := assert.New(t)

	input := "3,000 x 10^9"
	expected := Tokens{
		NewToken(tokenNumber, 0, "3,000 x 10^9"),
	}
	actual := NewLexer(input).Drain()
	a.Equal(expected, actual)
}
