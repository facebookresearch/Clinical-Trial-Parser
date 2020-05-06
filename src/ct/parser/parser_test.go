// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var parser *Parser

func init() {
	parser = NewParser()
}

func TestOneNumericalVariableParser(t *testing.T) {
	a := assert.New(t)

	input := "a1c greater than or equal to 5.0%."
	expected := List{
		Items{
			NewItem(itemVariable, "a1c"),
			NewItem(itemComparison, "≥"),
			NewItem(itemNumber, "5.0"),
			NewItem(itemUnit, "%"),
			NewItem(itemPunctuation, "."),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestToRangeParser(t *testing.T) {
	a := assert.New(t)

	input := "hba1c >6.0% to <11.5%"
	expected := List{
		Items{
			NewItem(itemVariable, "a1c"),
			NewItem(itemComparison, ">"),
			NewItem(itemNumber, "6.0"),
			NewItem(itemUnit, "%"),
			NewItem(itemComparison, "<"),
			NewItem(itemNumber, "11.5"),
			NewItem(itemUnit, "%"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestOneCategoricalVariableParser(t *testing.T) {
	a := assert.New(t)

	input := "eastern cooperative oncology group (ecog) performance status  0-3"
	expected := List{
		Items{
			NewItem(itemVariable, "ecog"),
			NewItem(itemNumber, "0"),
			NewItem(itemRange, "-"),
			NewItem(itemNumber, "3"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)

	input = "eastern cooperative oncology group (ecog) performance status of 0, 1, or 2."
	expected = List{
		Items{
			NewItem(itemVariable, "ecog"),
			NewItem(itemNumber, "0"),
			NewItem(itemNumber, "1"),
			NewItem(itemOr, "or"),
			NewItem(itemNumber, "2"),
			NewItem(itemPunctuation, "."),
		},
	}
	actual = parser.Parse(input)
	a.Equal(expected, actual)
}

func TestRomanNumeralsParser(t *testing.T) {
	a := assert.New(t)

	input := "nyha class iii or iv."
	expected := List{
		Items{
			NewItem(itemVariable, "nyha"),
			NewItem(itemNumber, "iii"),
			NewItem(itemOr, "or"),
			NewItem(itemNumber, "iv"),
			NewItem(itemPunctuation, "."),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestReverseOrderParser(t *testing.T) {
	a := assert.New(t)

	input := "patient < 30.0 kilograms (kg) in weight"
	expected := List{
		Items{
			NewItem(itemComparison, "<"),
			NewItem(itemNumber, "30.0"),
			NewItem(itemUnit, "kg"),
			NewItem(itemVariable, "weight"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestMisspelledVariableNameParser(t *testing.T) {
	a := assert.New(t)

	input := "weigh less than 110 pounds"
	expected := List{
		Items{
			NewItem(itemVariable, "weight"),
			NewItem(itemComparison, "<"),
			NewItem(itemNumber, "110"),
			NewItem(itemUnit, "lb"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestMissingVariableNameParser(t *testing.T) {
	a := assert.New(t)

	input := "patients less than 31 pounds"
	expected := List{
		Items{
			NewItem(itemComparison, "<"),
			NewItem(itemNumber, "31"),
			NewItem(itemUnit, "lb"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestTwoVariablesParser(t *testing.T) {
	a := assert.New(t)

	input := "weight > 50 kg and a1c <= 9.5%"
	expected := List{
		Items{
			NewItem(itemVariable, "weight"),
			NewItem(itemComparison, ">"),
			NewItem(itemNumber, "50"),
			NewItem(itemUnit, "kg"),
			NewItem(itemAnd, "and"),
			NewItem(itemVariable, "a1c"),
			NewItem(itemComparison, "≤"),
			NewItem(itemNumber, "9.5"),
			NewItem(itemUnit, "%"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestTwoVariablesWithMoreStuffParser(t *testing.T) {
	a := assert.New(t)

	input := "weight ≥ 50 kg; with a body mass index ≤ 45 kg/m²"
	expected := List{
		Items{
			NewItem(itemVariable, "weight"),
			NewItem(itemComparison, "≥"),
			NewItem(itemNumber, "50"),
			NewItem(itemUnit, "kg"),
			NewItem(itemPunctuation, ";"),
			NewItem(itemVariable, "bmi"),
			NewItem(itemComparison, "≤"),
			NewItem(itemNumber, "45"),
			NewItem(itemUnit, "kg/m2"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestRangeWithToParser(t *testing.T) {
	a := assert.New(t)

	input := "body mass index of 18.0 to 31.0 kg/m^2, and a total body weight >50 kg"
	expected := List{
		Items{
			NewItem(itemVariable, "bmi"),
			NewItem(itemNumber, "18.0"),
			NewItem(itemRange, "to"),
			NewItem(itemNumber, "31.0"),
			NewItem(itemUnit, "kg/m2"),
			NewItem(itemAnd, "and"),
			NewItem(itemVariable, "weight"),
			NewItem(itemComparison, ">"),
			NewItem(itemNumber, "50"),
			NewItem(itemUnit, "kg"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestParenthesisParser(t *testing.T) {
	a := assert.New(t)

	input := "eastern cooperative oncology group (ecog) performance status =< 2 (karnofsky >= 50%)"
	expected := List{
		Items{
			NewItem(itemVariable, "karnofsky_score"),
			NewItem(itemComparison, "≥"),
			NewItem(itemNumber, "50"),
			NewItem(itemUnit, "%"),
		},
		Items{
			NewItem(itemVariable, "ecog"),
			NewItem(itemComparison, "≤"),
			NewItem(itemNumber, "2"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestTwoVariablesInParenthesisParser(t *testing.T) {
	a := assert.New(t)

	input := "uncontrolled hypertension (systolic blood pressure (sbp) >140mmhg, diastolic blood pressure (dbp) >90mmhg)"
	expected := List{
		Items{
			NewItem(itemVariable, "sbp"),
			NewItem(itemComparison, ">"),
			NewItem(itemNumber, "140"),
			NewItem(itemUnit, "mmhg"),
			NewItem(itemVariable, "dbp"),
			NewItem(itemComparison, ">"),
			NewItem(itemNumber, "90"),
			NewItem(itemUnit, "mmhg"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestOneVariableWithTwoValuesParser(t *testing.T) {
	a := assert.New(t)

	input := "uncontrolled hypertension, defined as blood pressure (bp) > 140/90"
	expected := List{
		Items{
			NewItem(itemVariable, "sbp/dbp"),
			NewItem(itemComparison, ">"),
			NewItem(itemNumber, "140/90"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestRatioVariableParser(t *testing.T) {
	a := assert.New(t)

	input := "aspartate aminotransferase (ast)/alanine aminotransferase (alt) ≤ 2.0 x upper limits of normal"
	expected := List{
		Items{
			NewItem(itemVariable, "ast"),
			NewItem(itemSlash, "/"),
			NewItem(itemVariable, "alt"),
			NewItem(itemComparison, "≤"),
			NewItem(itemNumber, "2.0"),
			NewItem(itemUnit, "uln"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestAtLeastParser(t *testing.T) {
	a := assert.New(t)

	input := "life expectancy of at least 2 years"
	expected := List{
		Items{
			NewItem(itemVariable, "life_expectancy"),
			NewItem(itemComparison, "≥"),
			NewItem(itemNumber, "2"),
			NewItem(itemUnit, "year"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}

func TestCountParser(t *testing.T) {
	a := assert.New(t)

	input := "wbc > 3,000/ul"
	expected := List{
		Items{
			NewItem(itemVariable, "wbc"),
			NewItem(itemComparison, ">"),
			NewItem(itemNumber, "3,000"),
			NewItem(itemUnit, "cells/ul"),
		},
	}
	actual := parser.Parse(input)
	a.Equal(expected, actual)
}
